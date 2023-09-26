// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package confighttp

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configcompression"
)

func TestCompressor(t *testing.T) {
	testData := []byte("test")
	testCases := []struct{
		compressionType configcompression.CompressionType
		expected []byte
	}{
		{
			compressionType: configcompression.Zstd,
			expected: compressZstd(t, testData).Bytes(),
		},
	}

	for _, tC := range testCases {
		t.Run(string(tC.compressionType), func(t *testing.T) {
			compressor, err := newCompressor(tC.compressionType)
			require.NoError(t, err)
			readCloser := io.NopCloser(bytes.NewBuffer(testData))
			actual := bytes.NewBuffer(nil)
			err = compressor.compress(actual, readCloser)
			require.NoError(t, err)
			require.Equal(t, tC.expected, actual.Bytes())
		})
	}
}

func BenchmarkZstdCompressor(b *testing.B) {
	testData  := make([]byte, 100000)
	_, err := rand.Read(testData)
	require.NoError(b, err)
	expected := compressZstd(b, testData)

	compressor, err := newCompressor(configcompression.Zstd)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		readCloser := io.NopCloser(bytes.NewBuffer(testData))
		actual := bytes.NewBuffer(nil)
		b.StartTimer()
		err = compressor.compress(actual, readCloser)
		//time.Sleep(100*time.Millisecond)
		b.StopTimer()
		require.NoError(b, err)
		require.Equal(b, expected.Bytes(), actual.Bytes())
	}
}

func BenchmarkGzipCompressor(b *testing.B) {
	const (
		messageBlock       = "This is an example log"
	)

	message := strings.Repeat(messageBlock, 100000)
	testData := []byte(message)
	expected := compressGzip(b, testData)

	compressor, err := newCompressor(configcompression.Gzip)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		readCloser := io.NopCloser(bytes.NewBuffer(testData))
		actual := bytes.NewBuffer(nil)
		b.StartTimer()
		err = compressor.compress(actual, readCloser)
		b.StopTimer()
		require.NoError(b, err)
		require.Equal(b, expected.Bytes(), actual.Bytes())
	}
}

func TestZstdCompressor2(t *testing.T) {
	// around 200kb, which is a reasonable estimate for the strings we'd be compressing in practice
	testData  := make([]byte, 10000)
	_, err := rand.Read(testData)
	require.NoError(t, err)
	expected := compressZstd(t, testData)

	compressor, err := newCompressor(configcompression.Zstd)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		readCloser := io.NopCloser(bytes.NewBuffer(testData))
		actual := bytes.NewBuffer(nil)
		err = compressor.compress(actual, readCloser)
		time.Sleep(1000*time.Millisecond)
		require.NoError(t, err)
		require.Equal(t, expected.Bytes(), actual.Bytes())
	}
}
