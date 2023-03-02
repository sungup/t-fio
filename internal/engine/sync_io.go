package engine

import "os"

type SyncIO struct {
	fp *os.File
}

func (f *SyncIO) ReadAt(p []byte, offset int64, callback Callback) (err error) {
	n, e := f.fp.ReadAt(p, offset)

	callback(n, e)

	return e
}

func (f *SyncIO) WriteAt(p []byte, offset int64, callback Callback) (err error) {
	n, e := f.fp.WriteAt(p, offset)

	callback(n, e)

	return e
}

func (f *SyncIO) Close() (err error) {
	return f.fp.Close()
}
