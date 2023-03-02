package sys

import "io"

// File is an abstracted file IO interface to implement IO engine.
type File interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
}
