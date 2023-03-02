package engine

import "io"

type Callback func(n int, err error)

type Engine interface {
	ReadAt(p []byte, off int64, callback Callback) (err error)
	WriteAt(p []byte, off int64, callback Callback) (err error)
	io.Closer
}
