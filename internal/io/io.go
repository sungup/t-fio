package io

import "os"

type IO struct {
	id      int64  // Identified number
	offset  int64  // byte unit io position
	buffer  []byte // it should be aligned block
	next    *IO    // linked list item to the next transaction
	success bool

	issue func(fp *os.File, offset int64, buf []byte, callback func(success bool)) (err error)

	ch chan<- *IO
}

func (io *IO) Issue(fp *os.File) (next *IO, err error) {
	if err = io.issue(fp, io.offset, io.buffer, io.Callback); err == nil {
		next = io.next
	}

	return next, err
}

func (io *IO) Callback(success bool) {
	io.success = success

	io.ch <- io
}
