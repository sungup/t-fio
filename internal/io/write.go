package io

import (
	"fmt"
	"github.com/sungup/t-fio/pkg/sys"
)

func Write(fp sys.File, offset int64, buf []byte, callback func(success bool)) error {
	if fp == nil {
		return fmt.Errorf("issue target file should not be nil")
	}

	// Write async IO
	go func() {
		// TODO add write error handler for the write fail
		_, _ = fp.WriteAt(buf, offset)
	}()

	// always ack that this writing transaction has been success
	callback(true)

	return nil
}
