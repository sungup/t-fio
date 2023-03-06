//go:build linux
// +build linux

package engine

import (
	"context"
	"github.com/iceber/iouring-go"
	"os"
)

type IOURing struct {
	fp    *os.File
	uring *iouring.IOURing

	ch chan iouring.Result

	handlerCount int
}

func (f *IOURing) ReadAt(p []byte, offset int64, callback Callback) (err error) {
	prepReq := iouring.Pread(int(f.fp.Fd()), p, uint64(offset)).WithCallback(func(result iouring.Result) error {
		n, e := result.ReturnInt()
		callback(n, e)
		return e
	})

	_, err = f.uring.SubmitRequest(prepReq, f.ch)

	return err
}

func (f *IOURing) WriteAt(p []byte, offset int64, callback Callback) (err error) {
	prepReq := iouring.Pwrite(int(f.fp.Fd()), p, uint64(offset)).WithCallback(func(result iouring.Result) error {
		n, e := result.ReturnInt()
		callback(n, e)
		return e
	})

	_, err = f.uring.SubmitRequest(prepReq, f.ch)

	return err
}

func (f *IOURing) GetIOFunc(ioType IOType) (io DoIO, err error) {
	return getIOFunc(f, ioType)
}

func (f *IOURing) Close() (err error) {
	if err = f.uring.Close(); err == nil {
		err = f.fp.Close()
	}

	return err
}

func (f *IOURing) Run(ctx context.Context) {
	for i := 0; i < f.handlerCount; i++ {
		// start handling routines
		go func(ch <-chan iouring.Result, ctx context.Context) {
			for {
				select {
				case result := <-ch:
					_ = result.Callback()
				case <-ctx.Done():
					return
				}
			}
		}(f.ch, ctx)
	}
}
