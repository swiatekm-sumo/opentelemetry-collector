// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package configstorage // import "go.opentelemetry.io/collector/config/configstorage"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/experimental/storage"
)

var (
	errNoStorageClient    = errors.New("no storage client extension found")
	errWrongExtensionType = errors.New("requested extension is not a storage extension")
)

// NewDefaultBackOffConfig returns the default settings for StorageConfig.
func NewDefaultStorageConfig() StorageConfig {
	return StorageConfig{}
}

// StorageConfig defines configuration for obtaining a storage client.
type StorageConfig struct {
	// StorageID if not empty, uses the component specified  as a storage extension for the component
	StorageID *component.ID `mapstructure:"storage"`
}

func (sc *StorageConfig) Validate() error {
	return nil
}

func (sc *StorageConfig) ToStorageClient(
	ctx context.Context,
	host component.Host,
	componentID component.ID,
	componentKind component.Kind,
	signal component.DataType) (storage.Client, error) {

	if sc.StorageID == nil {
		return nil, errNoStorageClient
	}

	ext, found := host.GetExtensions()[*sc.StorageID]
	if !found {
		return nil, errNoStorageClient
	}

	storageExt, ok := ext.(storage.Extension)
	if !ok {
		return nil, errWrongExtensionType
	}

	return storageExt.GetClient(ctx, componentKind, componentID, signal.String())
}
