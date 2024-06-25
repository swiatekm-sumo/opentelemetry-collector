// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package configstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultStorageSettings(t *testing.T) {
	cfg := NewDefaultStorageConfig()
	assert.NoError(t, cfg.Validate())
	assert.Equal(t, StorageConfig{}, cfg)
}
