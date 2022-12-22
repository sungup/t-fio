package io

import "os"

func AsyncWrite(fp *os.File, offset int64, buf []byte, callback func(success bool)) error {
	// Read Sync IO
	go func() {
		// TODO add write error handler for the write fail
		_, _ = fp.WriteAt(buf, offset)
	}()

	// always ack that this writing transaction has been success
	callback(true)

	return nil
}
