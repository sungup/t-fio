package io

import (
	"github.com/sungup/t-fio/pkg/measure"
	"os"
	"sync"
	"time"
)

type IO struct {
	jobId  int64  // Identified number
	offset int64  // byte unit io position
	buffer []byte // it should be aligned block

	issue   func(fp *os.File, offset int64, buf []byte, callback func(success bool)) (err error)
	latency func() time.Duration

	wait *sync.WaitGroup
}

func (io *IO) Issue(fp *os.File, wait *sync.WaitGroup) error {
	io.wait = wait
	io.latency = measure.LatencyMeasureStart()
	return io.issue(fp, io.offset, io.buffer, io.Callback)
}

func (io *IO) Callback(success bool) {
	defer func() {
		// TODO call stat collector function
		io.latency()
	}()

	io.wait.Done()
}

func NewIO(ioType Type, jobId, offset int64, buffer []byte) *IO {
	return &IO{jobId: jobId, offset: offset, buffer: buffer, issue: ioType}
}
