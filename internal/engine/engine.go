package engine

import (
	"fmt"
	"io"
)

type Callback func(n int, err error)
type DoIO func(p []byte, offset int64, callback Callback) (err error)

type IOType int

const (
	Read = IOType(iota)
	Write
)

type Engine interface {
	ReadAt(p []byte, offset int64, callback Callback) (err error)
	WriteAt(p []byte, offset int64, callback Callback) (err error)
	GetIOFunc(ioType IOType) (io DoIO, err error)
	io.Closer
}

func getIOFunc(engine Engine, ioType IOType) (io DoIO, err error) {
	switch ioType {
	case Read:
		io = engine.ReadAt
	case Write:
		io = engine.WriteAt
	default:
		err = fmt.Errorf("unsupported IO type")
	}

	return
}
