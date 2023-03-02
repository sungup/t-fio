//go:build linux
// +build linux

package engine

import (
	"context"
	"github.com/iceber/iouring-go"
	"os"
)

type IOUring struct {
	fd    *os.File
	uring *iouring.IOURing

	ch chan iouring.Result

	handlerCount int
}

func (f *IOUring) ReadAt(p []byte, offset int64, callback Callback) (err error) {
	prepReq := iouring.Pread(int(f.fd.Fd()), p, uint64(offset)).WithCallback(func(result iouring.Result) error {
		n, e := result.ReturnInt()
		callback(n, e)
		return e
	})

	_, err = f.uring.SubmitRequest(prepReq, f.ch)

	return
}

func (f *IOUring) WriteAt(p []byte, offset int64, callback Callback) (err error) {
	prepReq := iouring.Pwrite(int(f.fd.Fd()), p, uint64(offset)).WithCallback(func(result iouring.Result) error {
		n, e := result.ReturnInt()
		callback(n, e)
		return e
	})

	_, err = f.uring.SubmitRequest(prepReq, f.ch)

	return
}

func (f *IOUring) Close() (err error) {
	if err = f.uring.Close(); err == nil {
		err = f.fd.Close()
	}

	return
}

func (f *IOUring) runResultHandler(ctx context.Context) {
	for i := 0; i < f.handlerCount; i++ {
		// start handling routines
		go func(ch <-chan iouring.Result) {
			for {
				select {
				case result := <-ch:
					_ = result.Callback()
				case <-ctx.Done():
					return
				}
			}
		}(f.ch)
	}
}
