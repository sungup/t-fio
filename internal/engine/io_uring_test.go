//go:build linux
// +build linux

package engine

import (
	"context"
	"github.com/iceber/iouring-go"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
	"time"
)

const tcIourQd = 8

func tcInitIOURing() (iour *IOURing, closer func()) {
	ctx, ctxCloser := context.WithCancel(context.Background())
	ch := make(chan iouring.Result, tcIourQd)
	uring, _ := iouring.New(tcIourQd)

	iour = &IOURing{
		fp:           nil,
		uring:        uring,
		ch:           ch,
		ctx:          ctx,
		handlerCount: tcIourQd,
	}

	for i := 0; i < tcIourQd; i++ {
		go func(ch <-chan iouring.Result, ctx context.Context) {
			for {
				select {
				case result := <-ch:
					_ = result.Callback()
				case <-ctx.Done():
					return
				}
			}
		}(ch, ctx)
	}

	return iour, func() {
		ctxCloser()
		_ = uring.Close()
	}
}

func TestIOUring_ReadAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested, closer := tcInitIOURing()
	defer closer()

	tcCounter.Add(1)
	assert.NoError(t, tested.ReadAt(nil, 0, tcFailedCB))
	assert.NotZero(t, tcCounter.Len())

	// wait until all thread completed
	tcCounter.Wait()

	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestIOUring_ReadAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested.fp = tcFile

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		var testedBuffer [test.BufferSz]byte

		tcCounter.Add(1)
		assert.NoError(t, tested.ReadAt(testedBuffer[:], tcOffset, tcSuccessCB))
		assert.NotZero(t, tcCounter.Len())
	}

	// wait until all thread completed
	tcCounter.Wait()
}

func TestIOUring_WriteAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	writtenBuffer := make([]byte, test.BufferSz)
	expectedBuffer := make([]byte, test.BufferSz)
	test.FillBuffer(expectedBuffer, time.Now().UnixNano())

	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested, closer := tcInitIOURing()
	defer closer()

	// fail test
	tcCounter.Add(1)
	assert.NoError(t, tested.WriteAt(nil, 0, tcFailedCB))
	assert.NotZero(t, tcCounter.Len())

	// wait until all thread completed
	tcCounter.Wait()
	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestIOUring_WriteAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested.fp = tcFile

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		// check write is success
		tcCounter.Add(1)
		assert.NoError(t, tested.WriteAt(expectedBuffer, tcOffset, tcSuccessCB))
		assert.NotZero(t, tcCounter.Len())
	}

	// wait until all thread completed
	tcCounter.Wait()

	// check written data after all io completed
	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		_, err = tcFile.ReadAt(writtenBuffer, tcOffset)
		assert.NoError(t, err)
		assert.Equal(t, expectedBuffer, writtenBuffer)
		assert.NotEqual(t, test.Buffer, writtenBuffer)
	}

}

func TestIOUring_Close(t *testing.T) {
	tcFile, tcCloser, err := test.OpenTCFile("TestIOUring_Close", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	// NoError
	tested, _ := tcInitIOURing()
	tested.fp = tcFile
	assert.NoError(t, tested.Close())
	assert.True(t, tested.uring.IsClosed())

	// error with already closed file, and there is no way normally raise iouring.close error
	tested, _ = tcInitIOURing()
	tested.fp = tcFile
	assert.False(t, tested.uring.IsClosed())
	assert.Error(t, tested.Close())
	assert.True(t, tested.uring.IsClosed())
}
