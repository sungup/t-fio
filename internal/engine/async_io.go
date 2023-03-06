package engine

import (
	"os"
)

type AsyncIO struct {
	fp *os.File
}

func (f *AsyncIO) ReadAt(p []byte, offset int64, callback Callback) (err error) {
	go func() {
		n, e := f.fp.ReadAt(p, offset)

		callback(n, e)
	}()

	return nil
}

func (f *AsyncIO) WriteAt(p []byte, offset int64, callback Callback) (err error) {
	go func() {
		n, e := f.fp.WriteAt(p, offset)

		callback(n, e)
	}()

	return nil
}

func (f *AsyncIO) GetIOFunc(ioType IOType) (io DoIO, err error) {
	return getIOFunc(f, ioType)
}

func (f *AsyncIO) Close() (err error) {
	return f.fp.Close()
}
