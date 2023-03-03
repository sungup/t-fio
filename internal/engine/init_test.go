package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
	"time"
)

const (
	tcFileSz = test.BufferSz * 1024
	tcDelay  = time.Microsecond * 50
)

func makeFailedCallback(t *testing.T, cnt *test.AtomicCounter) Callback {
	return func(n int, err error) {
		assert.LessOrEqual(t, n, 0)
		assert.Error(t, err)
		time.Sleep(tcDelay) // add delay time to avoid too fast callback
		cnt.Done()
	}
}

func makeSuccessCallback(t *testing.T, cnt *test.AtomicCounter) Callback {
	return func(n int, err error) {
		assert.Equal(t, test.BufferSz, n)
		assert.NoError(t, err)
		time.Sleep(tcDelay) // add delay time to avoid too fast callback
		cnt.Done()
	}
}
