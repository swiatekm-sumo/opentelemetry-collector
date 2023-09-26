// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package confighttp // import "go.opentelemetry.io/collector/config/confighttp"

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"runtime"

	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"

	"go.opentelemetry.io/collector/config/configcompression"
)

type writeCloserReset interface {
	io.WriteCloser
	Reset(w io.Writer)
}

var (
	_          writeCloserReset = (*gzip.Writer)(nil)
	gZipPool   *compressor
	_          writeCloserReset = (*snappy.Writer)(nil)
	snappyPool *compressor
	_          writeCloserReset = (*zstd.Encoder)(nil)
	zStdPool   *compressor
	_          writeCloserReset = (*zlib.Writer)(nil)
	zLibPool   *compressor
)

func init() {
	maxPoolSize := runtime.NumCPU()
	gZipPool = &compressor{newEncoderPool(maxPoolSize, func() writeCloserReset { return gzip.NewWriter(nil) })}
	snappyPool = &compressor{newEncoderPool(maxPoolSize, func() writeCloserReset { return snappy.NewBufferedWriter(nil) })}
	zStdPool = &compressor{newEncoderPool(maxPoolSize, func() writeCloserReset { zw, _ := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1)); return zw })}
	zLibPool = &compressor{newEncoderPool(maxPoolSize, func() writeCloserReset { return zlib.NewWriter(nil) })}

}

type compressor struct {
	pool *encoderPool
}

// writerFactory defines writer field in CompressRoundTripper.
// The validity of input is already checked when NewCompressRoundTripper was called in confighttp,
func newCompressor(compressionType configcompression.CompressionType) (*compressor, error) {
	switch compressionType {
	case configcompression.Gzip:
		return gZipPool, nil
	case configcompression.Snappy:
		return snappyPool, nil
	case configcompression.Zstd:
		return zStdPool, nil
	case configcompression.Zlib, configcompression.Deflate:
		return zLibPool, nil
	}
	return nil, errors.New("unsupported compression type, ")
}

func (p *compressor) compress(buf *bytes.Buffer, body io.ReadCloser) error {
	writer := p.pool.Get()
	defer p.pool.Put(writer)
	writer.Reset(buf)

	if body != nil {
		_, copyErr := io.Copy(writer, body)
		closeErr := body.Close()

		if copyErr != nil {
			return copyErr
		}

		if closeErr != nil {
			return closeErr
		}
	}

	return writer.Close()
}

type encoderPool struct {
	New func() writeCloserReset
	pool chan writeCloserReset
}

func newEncoderPool(maxSize int, new func() writeCloserReset) *encoderPool {
	return &encoderPool{
		New: new,
		pool: make(chan writeCloserReset, maxSize),
	}
}

func (p *encoderPool) Get() writeCloserReset {
	select {
	case encoder := <-p.pool:
		return encoder
	default:
		return p.New()
	}
}

func (p *encoderPool) Put(encoder writeCloserReset) {
	select {
	case p.pool <- encoder:
	default:
	}
}