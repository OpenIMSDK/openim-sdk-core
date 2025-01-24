package interaction

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"sync"

	"github.com/openimsdk/tools/errs"
)

var (
	gzipWriterPool = sync.Pool{New: func() any { return gzip.NewWriter(nil) }}
	gzipReaderPool = sync.Pool{New: func() any { return new(gzip.Reader) }}
)

type Compressor interface {
	Compress(rawData []byte) ([]byte, error)
	CompressWithPool(rawData []byte) ([]byte, error)
	DeCompress(compressedData []byte) ([]byte, error)
	DecompressWithPool(compressedData []byte) ([]byte, error)
}

type GzipCompressor struct {
	compressProtocol string
}

func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{compressProtocol: "gzip"}
}

func (g *GzipCompressor) Compress(rawData []byte) ([]byte, error) {
	gzipBuffer := bytes.Buffer{}
	gz := gzip.NewWriter(&gzipBuffer)
	if _, err := gz.Write(rawData); err != nil {
		return nil, errs.WrapMsg(err, "")
	}
	if err := gz.Close(); err != nil {
		return nil, errs.WrapMsg(err, "")
	}
	return gzipBuffer.Bytes(), nil
}

func (g *GzipCompressor) CompressWithPool(rawData []byte) ([]byte, error) {
	gz := gzipWriterPool.Get().(*gzip.Writer)
	defer gzipWriterPool.Put(gz)

	gzipBuffer := bytes.Buffer{}
	gz.Reset(&gzipBuffer)

	if _, err := gz.Write(rawData); err != nil {
		return nil, errs.WrapMsg(err, "")
	}
	if err := gz.Close(); err != nil {
		return nil, errs.WrapMsg(err, "")
	}
	return gzipBuffer.Bytes(), nil
}

func (g *GzipCompressor) DeCompress(compressedData []byte) ([]byte, error) {
	buff := bytes.NewBuffer(compressedData)
	reader, err := gzip.NewReader(buff)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewReader failed")
	}
	compressedData, err = io.ReadAll(reader)
	if err != nil {
		return nil, errs.WrapMsg(err, "ReadAll failed")
	}
	_ = reader.Close()
	return compressedData, nil
}

func (g *GzipCompressor) DecompressWithPool(compressedData []byte) ([]byte, error) {
	reader := gzipReaderPool.Get().(*gzip.Reader)
	if reader == nil {
		return nil, errors.New("NewReader failed")
	}
	defer gzipReaderPool.Put(reader)

	err := reader.Reset(bytes.NewReader(compressedData))
	if err != nil {
		return nil, errs.WrapMsg(err, "NewReader failed")
	}

	compressedData, err = io.ReadAll(reader)
	if err != nil {
		return nil, errs.WrapMsg(err, "ReadAll failed")
	}
	_ = reader.Close()
	return compressedData, nil
}
