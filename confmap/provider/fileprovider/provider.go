// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileprovider // import "go.opentelemetry.io/collector/confmap/provider/fileprovider"

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/internal"
)

const schemeName = "file"

type provider struct{}

// New returns a new confmap.Provider that reads the configuration from a file.
//
// This Provider supports "file" scheme, and can be called with a "uri" that follows:
//
//	file-uri		= "file:" local-path
//	local-path		= [ drive-letter ] file-path
//	drive-letter	= ALPHA ":"
//
// The "file-path" can be relative or absolute, and it can be any OS supported format.
//
// Examples:
// `file:path/to/file` - relative path (unix, windows)
// `file:/path/to/file` - absolute path (unix, windows)
// `file:c:/path/to/file` - absolute path including drive-letter (windows)
// `file:c:\path\to\file` - absolute path including drive-letter (windows)
func New() confmap.Provider {
	return &provider{}
}

func (fmp *provider) Retrieve(ctx context.Context, uri string, watcher confmap.WatcherFunc) (confmap.Retrieved, error) {
	if !strings.HasPrefix(uri, schemeName+":") {
		return confmap.Retrieved{}, fmt.Errorf("%q uri is not supported by %q provider", uri, schemeName)
	}

	// Clean the path before using it.
	cleanedPath := filepath.Clean(uri[len(schemeName)+1:])
	content, err := ioutil.ReadFile(cleanedPath)
	if err != nil {
		return confmap.Retrieved{}, fmt.Errorf("unable to read the file %v: %w", uri, err)
	}

	opts := []confmap.RetrievedOption{}
	if watcher != nil {
		closeFunc, err := watchConfigFile(ctx, cleanedPath, watcher)
		if err != nil {
			return confmap.Retrieved{}, fmt.Errorf("unable to start file watcher %v: %w", uri, err)
		}
		opts = append(opts, confmap.WithRetrievedClose(closeFunc))
	}

	return internal.NewRetrievedFromYAML(content, opts...)
}

func (*provider) Scheme() string {
	return schemeName
}

func (fmp *provider) Shutdown(context.Context) error {
	return nil
}

// watchConfigFile starts a file watcher for the provided path, which monitors it for changes and converts
// them to change events for the provided config watcher
// returns a function to close the watcher, or an error
func watchConfigFile(ctx context.Context, path string, configWatcher confmap.WatcherFunc) (confmap.CloseFunc, error) {
	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to create file watcher for %v: %w", path, err)
	}

	if err = fileWatcher.Add(path); err != nil {
		return nil, fmt.Errorf("unable to watch the file %v: %w", path, err)
	}

	go startFileEventLoop(ctx, fileWatcher, configWatcher)

	closeFunc := func(_ context.Context) error {
		return fileWatcher.Close()
	}

	return closeFunc, nil
}

// startFileEventLoop starts a loop which processes filesystem events from the fileWatcher, converts them
// into ChangeEvents, and calls the configWatcher function on them
// returns when the fileWatcher is closed, or when the context is cancelled
func startFileEventLoop(ctx context.Context, fileWatcher *fsnotify.Watcher, configWatcher confmap.WatcherFunc) {
	for { // this loop will exit once the fileWatcher is closed or the context is cancelled
		select {
		case fsEvent, ok := <-fileWatcher.Events:
			if !ok {
				return
			}
			changeEvent := fsEventToChangeEvent(fsEvent)
			if changeEvent != nil {
				configWatcher(changeEvent)
			}
		case err, ok := <-fileWatcher.Errors:
			if !ok {
				return
			}
			changeEvent := confmap.ChangeEvent{Error: err}
			configWatcher(&changeEvent)
		case <-ctx.Done():
			return
		}
	}
}

// fsEventToChangeEvent converts file system events from fsnotify to ChangeEvents for the config provider
// the logic is that we emit a normal change event for writes and creates, errors for renames and removes
// and we ignore permission changes
func fsEventToChangeEvent(fsEvent fsnotify.Event) *confmap.ChangeEvent {
	var err error
	switch fsEvent.Op {
	case fsnotify.Write, fsnotify.Create:
		err = nil
	case fsnotify.Rename:
		err = fmt.Errorf("config file moved: %v", fsEvent.Name)
	case fsnotify.Remove:
		err = fmt.Errorf("config file removed: %v", fsEvent.Name)
	case fsnotify.Chmod: // no event, we don't really care
		return nil
	}
	return &confmap.ChangeEvent{Error: err}
}
