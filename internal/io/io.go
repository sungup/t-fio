package io

import (
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"sync"
)

type IO struct {
	jobId  int64            // Identified number
	offset int64            // byte unit io position
	buffer *bytebuf.ByteBuf // it should be aligned block

	issue engine.DoIO
}

func (io *IO) Issue(wait *sync.WaitGroup) error {
	latency := measure.LatencyMeasureStart()

	return io.issue(io.buffer.Buffer(), io.offset, func(n int, err error) {
		defer func() {
			// TODO call stat collector function
			_ = latency()

			// return buffeer to memory pool to reuse
			io.buffer.Close()
		}()

		wait.Done()
	})
}

func New(doIO engine.DoIO, jobId, offset int64, buffer *bytebuf.ByteBuf) *IO {
	return &IO{jobId: jobId, offset: offset, buffer: buffer, issue: doIO}
}
