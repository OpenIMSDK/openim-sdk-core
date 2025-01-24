package file

import (
	"io"
)

func NewProgressReader(r io.Reader, fn func(current int64)) io.Reader {
	if r == nil || fn == nil {
		return r
	}
	return &Reader{
		r:  r,
		fn: fn,
	}
}

type Reader struct {
	r    io.Reader
	read int64
	fn   func(current int64)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == nil && n > 0 {
		r.read += int64(n)
		r.fn(r.read)
	}
	return n, err
}
