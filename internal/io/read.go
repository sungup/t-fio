package io

import "os"

func SyncRead(fp *os.File, offset int64, buf []byte, callback func(success bool)) error {
	// Read Sync IO
	_, err := fp.ReadAt(buf, offset)

	// run callback any time
	callback(err == nil)

	return err
}

func AsyncRead(fp *os.File, offset int64, buf []byte, callback func(success bool)) error {
	// Read Sync IO
	go func() {
		_ = SyncRead(fp, offset, buf, callback)
	}()

	return nil
}
