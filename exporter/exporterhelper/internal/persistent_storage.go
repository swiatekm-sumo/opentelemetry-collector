// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build enable_unstable
// +build enable_unstable

package internal // import "go.opentelemetry.io/collector/exporter/exporterhelper/internal"

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/extension/experimental/storage"
)

// persistentStorage provides an interface for request storage operations
type persistentStorage interface {
	// put appends the request to the storage
	put(req PersistentRequest) error
	// get returns the next available request; note that the channel is unbuffered
	get() <-chan PersistentRequest
	// size returns the current size of the persistent storage with items waiting for processing
	size() uint64
	// stop gracefully stops the storage
	stop()
}

// persistentContiguousStorage provides a persistent queue implementation backed by file storage extension
//
// Write index describes the position at which next item is going to be stored.
// Read index describes which item needs to be read next.
// When Write index = Read index, no elements are in the queue.
//
// The items currently dispatched by consumers are not deleted until the processing is finished.
// Their list is stored under a separate key.
//
//
//   ┌───────file extension-backed queue───────┐
//   │                                         │
//   │     ┌───┐     ┌───┐ ┌───┐ ┌───┐ ┌───┐   │
//   │ n+1 │ n │ ... │ 4 │ │ 3 │ │ 2 │ │ 1 │   │
//   │     └───┘     └───┘ └─x─┘ └─|─┘ └─x─┘   │
//   │                       x     |     x     │
//   └───────────────────────x─────|─────x─────┘
//      ▲              ▲     x     |     x
//      │              │     x     |     xxxx deleted
//      │              │     x     |
//    write          read    x     └── currently dispatched item
//    index          index   x
//                           xxxx deleted
//
type persistentContiguousStorage struct {
	logger      *zap.Logger
	queueName   string
	client      storage.Client
	unmarshaler RequestUnmarshaler

	putChan  chan struct{}
	stopChan chan struct{}
	stopOnce sync.Once
	capacity uint64

	reqChan chan PersistentRequest

	mu                       sync.Mutex
	readIndex                itemIndex
	writeIndex               itemIndex
	currentlyDispatchedItems []itemIndex

	itemsCount uint64
}

type itemIndex uint64

const (
	zapKey           = "key"
	zapQueueNameKey  = "queueName"
	zapErrorCount    = "errorCount"
	zapNumberOfItems = "numberOfItems"

	readIndexKey                = "ri"
	writeIndexKey               = "wi"
	currentlyDispatchedItemsKey = "di"
)

var (
	errMaxCapacityReached   = errors.New("max capacity reached")
	errValueNotSet          = errors.New("value not set")
	errKeyNotPresentInBatch = errors.New("key was not present in get batchStruct")
)

// newPersistentContiguousStorage creates a new file-storage extension backed queue;
// queueName parameter must be a unique value that identifies the queue.
// The queue needs to be initialized separately using initPersistentContiguousStorage.
func newPersistentContiguousStorage(ctx context.Context, queueName string, capacity uint64, logger *zap.Logger, client storage.Client, unmarshaler RequestUnmarshaler) *persistentContiguousStorage {
	pcs := &persistentContiguousStorage{
		logger:      logger,
		client:      client,
		queueName:   queueName,
		unmarshaler: unmarshaler,
		capacity:    capacity,
		putChan:     make(chan struct{}, capacity),
		reqChan:     make(chan PersistentRequest),
		stopChan:    make(chan struct{}),
	}

	initPersistentContiguousStorage(ctx, pcs)

	// Make sure the leftover requests are handled
	pcs.requeueNotDispatchedReqs(context.Background())

	// We start the loop first so in case there are more elements in the persistent storage than the capacity,
	// it does not get blocked on initialization

	go pcs.loop()

	// Make sure the communication channel is loaded up
	for i := uint64(0); i < pcs.size(); i++ {
		pcs.putChan <- struct{}{}
	}

	return pcs
}

func initPersistentContiguousStorage(ctx context.Context, pcs *persistentContiguousStorage) {
	var writeIndex itemIndex
	var readIndex itemIndex
	batch, err := newBatch(pcs).get(readIndexKey, writeIndexKey, currentlyDispatchedItemsKey).execute(ctx)

	if err == nil {
		readIndex, err = batch.getItemIndexResult(readIndexKey)
	}

	if err == nil {
		writeIndex, err = batch.getItemIndexResult(writeIndexKey)
	}

	if err != nil {
		if err == errValueNotSet {
			pcs.logger.Info("Initializing new persistent queue", zap.String(zapQueueNameKey, pcs.queueName))
			pcs.readIndex = 0
			pcs.writeIndex = 0
		} else {
			pcs.logger.Error("Failed getting read/write index, starting with new ones",
				zap.String(zapQueueNameKey, pcs.queueName),
				zap.Error(err))
			// if we have dispatched items, set the read and write indices to be greater than theirs
			// to avoid overwriting them
			var dispatchedItems []itemIndex
			dispatchedItems, err = batch.getItemIndexArrayResult(currentlyDispatchedItemsKey)
			if err != nil || len(dispatchedItems) == 0 {
				pcs.readIndex = 0
				pcs.writeIndex = 0
			} else {
				var maxItemIndex itemIndex = dispatchedItems[0]
				for _, value := range dispatchedItems {
					if maxItemIndex < value {
						maxItemIndex = value
					}
				}
				pcs.readIndex = maxItemIndex + 1
				pcs.writeIndex = maxItemIndex + 1
			}
		}

	} else {
		pcs.readIndex = readIndex
		pcs.writeIndex = writeIndex
	}

	atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))
}

// loop is the main loop that handles fetching items from the persistent buffer
func (pcs *persistentContiguousStorage) loop() {
	for {
		select {
		case <-pcs.stopChan:
			return
		case <-pcs.putChan:
			req, found := pcs.getNextItem(context.Background())
			if found {
				pcs.reqChan <- req
			}
		}
	}
}

// get returns the request channel that all the requests will be send on
func (pcs *persistentContiguousStorage) get() <-chan PersistentRequest {
	return pcs.reqChan
}

// size returns the number of currently available items, which were not picked by consumers yet
func (pcs *persistentContiguousStorage) size() uint64 {
	return atomic.LoadUint64(&pcs.itemsCount)
}

func (pcs *persistentContiguousStorage) stop() {
	pcs.logger.Debug("Stopping persistentContiguousStorage", zap.String(zapQueueNameKey, pcs.queueName))
	pcs.stopOnce.Do(func() {
		close(pcs.stopChan)
	})
}

// put marshals the request and puts it into the persistent queue
func (pcs *persistentContiguousStorage) put(req PersistentRequest) error {
	// Nil requests are ignored
	if req == nil {
		return nil
	}

	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	if pcs.size() >= pcs.capacity {
		pcs.logger.Warn("Maximum queue capacity reached", zap.String(zapQueueNameKey, pcs.queueName))
		return errMaxCapacityReached
	}

	itemKey := pcs.itemKey(pcs.writeIndex)
	pcs.writeIndex++
	atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))

	ctx := context.Background()
	_, err := newBatch(pcs).setItemIndex(writeIndexKey, pcs.writeIndex).setRequest(itemKey, req).execute(ctx)

	// Inform the loop that there's some data to process
	pcs.putChan <- struct{}{}

	return err
}

// getNextItem pulls the next available item from the persistent storage; if none is found, returns (nil, false)
func (pcs *persistentContiguousStorage) getNextItem(ctx context.Context) (PersistentRequest, bool) {
	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	if pcs.readIndex != pcs.writeIndex {
		index := pcs.readIndex
		// Increase here, so even if errors happen below, it always iterates
		pcs.readIndex++
		atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))

		pcs.updateReadIndex(ctx)
		pcs.itemDispatchingStart(ctx, index)

		var req PersistentRequest
		batch, err := newBatch(pcs).get(pcs.itemKey(index)).execute(ctx)
		if err == nil {
			req, err = batch.getRequestResult(pcs.itemKey(index))
		}

		if err != nil || req == nil {
			// We need to make sure that currently dispatched items list is cleaned
			pcs.itemDispatchingFinish(ctx, index)

			return nil, false
		}

		// If all went well so far, cleanup will be handled by callback
		req.SetOnProcessingFinished(func() {
			pcs.mu.Lock()
			defer pcs.mu.Unlock()
			pcs.itemDispatchingFinish(ctx, index)
		})
		return req, true
	}

	return nil, false
}

// requeueNotDispatchedReqs gets the items for which sending was not finished, cleans the storage
// and moves the items back to the queue. Items which cannot be read or unmarshalled are ignored.
func (pcs *persistentContiguousStorage) requeueNotDispatchedReqs(ctx context.Context) {
	var dispatchedItems []itemIndex

	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	pcs.logger.Debug("Checking if there are items left for dispatch by consumers", zap.String(zapQueueNameKey, pcs.queueName))
	batch, err := newBatch(pcs).get(currentlyDispatchedItemsKey).execute(ctx)
	if err == nil {
		dispatchedItems, err = batch.getItemIndexArrayResult(currentlyDispatchedItemsKey)
	}
	if err != nil {
		pcs.logger.Error("Could not fetch items left for dispatch by consumers", zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
		return
	}

	if len(dispatchedItems) > 0 {
		pcs.logger.Info("Fetching items left for dispatch by consumers",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Int(zapNumberOfItems, len(dispatchedItems)))
	} else {
		pcs.logger.Debug("No items left for dispatch by consumers")
	}

	for index, item := range dispatchedItems {
		itemKey := pcs.itemKey(item)
		retrieveBatch := newBatch(pcs).get(itemKey)
		_, retrieveErr := retrieveBatch.execute(ctx)
		if retrieveErr != nil {
			pcs.logger.Warn("Failed retrieving item left by consumers", zap.String(zapQueueNameKey, pcs.queueName), zap.Error(retrieveErr))
			_, _ = newBatch(pcs).setItemIndexArray(currentlyDispatchedItemsKey, dispatchedItems).execute(ctx)
			continue
		}
		req, err := retrieveBatch.getRequestResult(itemKey)
		// If error happened or item is nil, ignore it
		if err != nil {
			pcs.logger.Warn("Failed unmarshalling item",
				zap.String(zapQueueNameKey, pcs.queueName), zap.String(zapKey, itemKey), zap.Error(err))
			_, _ = newBatch(pcs).setItemIndexArray(currentlyDispatchedItemsKey, dispatchedItems).execute(ctx)
			continue
		} else if req == nil {
			pcs.logger.Debug("Item value could not be retrieved",
				zap.String(zapQueueNameKey, pcs.queueName), zap.String(zapKey, itemKey), zap.Error(err))
			continue
		}

		// this logic is similar to .put(), but different enough that it's difficult to make DRY
		requeueKey := pcs.itemKey(pcs.writeIndex)
		requeueBatch := newBatch(pcs)
		remainingItems := dispatchedItems[(index + 1):]
		requeueBatch.setItemIndexArray(currentlyDispatchedItemsKey, remainingItems).
			setItemIndex(writeIndexKey, pcs.writeIndex+1).
			setRequest(requeueKey, req).
			delete(itemKey)

		_, requeueErr := requeueBatch.execute(ctx)
		if requeueErr != nil {
			pcs.logger.Warn("Failed requeueing item",
				zap.String(zapQueueNameKey, pcs.queueName), zap.String(zapKey, itemKey), zap.Error(requeueErr))
			continue
		}

		// update local counters after successful requeue
		pcs.writeIndex++
		atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))
		pcs.putChan <- struct{}{}
	}
}

// itemDispatchingStart appends the item to the list of currently dispatched items
func (pcs *persistentContiguousStorage) itemDispatchingStart(ctx context.Context, index itemIndex) {
	pcs.currentlyDispatchedItems = append(pcs.currentlyDispatchedItems, index)
	_, err := newBatch(pcs).
		setItemIndexArray(currentlyDispatchedItemsKey, pcs.currentlyDispatchedItems).
		execute(ctx)
	if err != nil {
		pcs.logger.Debug("Failed updating currently dispatched items",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

// itemDispatchingFinish removes the item from the list of currently dispatched items and deletes it from the persistent queue
func (pcs *persistentContiguousStorage) itemDispatchingFinish(ctx context.Context, index itemIndex) {
	var updatedDispatchedItems []itemIndex
	for _, it := range pcs.currentlyDispatchedItems {
		if it != index {
			updatedDispatchedItems = append(updatedDispatchedItems, it)
		}
	}
	pcs.currentlyDispatchedItems = updatedDispatchedItems

	_, err := newBatch(pcs).
		setItemIndexArray(currentlyDispatchedItemsKey, pcs.currentlyDispatchedItems).
		delete(pcs.itemKey(index)).
		execute(ctx)
	if err != nil {
		pcs.logger.Debug("Failed updating currently dispatched items",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

func (pcs *persistentContiguousStorage) updateReadIndex(ctx context.Context) {
	_, err := newBatch(pcs).
		setItemIndex(readIndexKey, pcs.readIndex).
		execute(ctx)

	if err != nil {
		pcs.logger.Debug("Failed updating read index",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

func (pcs *persistentContiguousStorage) itemKey(index itemIndex) string {
	return strconv.FormatUint(uint64(index), 10)
}
