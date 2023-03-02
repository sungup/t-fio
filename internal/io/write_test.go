package io

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
	"time"
)

func TestAsyncWrite(t *testing.T) {
	a := assert.New(t)

	var (
		writtenBuffer  = make([]byte, test.BufferSz)
		expectedBuffer = make([]byte, test.BufferSz)
	)

	test.FillBuffer(expectedBuffer, time.Now().UnixNano())

	tcCallback := func(success bool) {
		a.True(success)
	}

	// Fail Test, but call back always true
	a.Error(Write(nil, 0, nil, tcCallback))

	// Success Test
	tcFile, tcCloser, err := test.OpenTCFile("TestAsyncWrite", tcFileSz)
	a.NoError(err)
	defer tcCloser()

	tcSuccessCallback := func(success bool) {
		a.True(success)
	}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		a.NoError(Write(tcFile, tcOffset, expectedBuffer, tcSuccessCallback))
	}

	// Flush all written data 100ms after
	time.Sleep(time.Millisecond * 100)
	a.NoError(tcFile.Sync())

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		_, err = tcFile.ReadAt(writtenBuffer, tcOffset)
		a.NoError(err)
		a.Equal(expectedBuffer, writtenBuffer)
		a.NotEqual(test.Buffer, writtenBuffer)
	}
}
