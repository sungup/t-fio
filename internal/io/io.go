package io

import (
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/pkg/sys"
	"sync"
	"time"
)

type IO struct {
	jobId  int64            // Identified number
	offset int64            // byte unit io position
	buffer *bytebuf.ByteBuf // it should be aligned block

	issue   func(fp sys.File, offset int64, buf []byte, callback func(success bool)) (err error)
	latency func() time.Duration

	wait *sync.WaitGroup
}

func (io *IO) Issue(fp sys.File, wait *sync.WaitGroup) error {
	io.wait = wait
	io.latency = measure.LatencyMeasureStart()
	return io.issue(fp, io.offset, io.buffer.Buffer(), io.Callback)
}

func (io *IO) Callback(success bool) {
	defer func() {
		// TODO call stat collector function
		io.latency()

		// return buffer to memory pool to reuse
		io.buffer.Close()
	}()

	io.wait.Done()
}

func New(ioType Type, jobId, offset int64, buffer *bytebuf.ByteBuf) *IO {
	return &IO{jobId: jobId, offset: offset, buffer: buffer, issue: ioType}
}
