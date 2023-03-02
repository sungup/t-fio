package io

import (
	"fmt"
	"github.com/sungup/t-fio/pkg/sys"
)

func SyncRead(fp sys.File, offset int64, buf []byte, callback func(success bool)) error {
	// Read Sync IO
	_, err := fp.ReadAt(buf, offset)

	// run callback any time
	callback(err == nil)

	return err
}

func AsyncRead(fp sys.File, offset int64, buf []byte, callback func(success bool)) error {
	if fp == nil {
		return fmt.Errorf("issue target file should not be nil")
	}

	// Read Sync IO
	go func() {
		_ = SyncRead(fp, offset, buf, callback)
	}()

	return nil
}
