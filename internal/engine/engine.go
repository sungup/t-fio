package engine

import "io"

type Callback func(n int, err error)

type Engine interface {
	ReadAt(p []byte, offset int64, callback Callback) (err error)
	WriteAt(p []byte, offset int64, callback Callback) (err error)
	io.Closer
}
