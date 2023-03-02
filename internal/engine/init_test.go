package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
)

const (
	tcFileSz = test.BufferSz * 1024
)

func makeFailedCallback(t *testing.T, cnt *test.AtomicCounter) Callback {
	return func(n int, err error) {
		assert.Zero(t, n)
		assert.EqualError(t, err, "invalid argument")
		cnt.Done()
	}
}

func makeSuccessCallback(t *testing.T, cnt *test.AtomicCounter) Callback {
	return func(n int, err error) {
		assert.Equal(t, test.BufferSz, n)
		assert.NoError(t, err)
		cnt.Done()
	}
}
