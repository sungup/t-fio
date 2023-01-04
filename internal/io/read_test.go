package io

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
)

func TestSyncRead(t *testing.T) {
	a := assert.New(t)

	// Fail Test
	tcFailedCallback := func(success bool) {
		a.False(success)
	}

	a.Error(SyncRead(nil, 0, nil, tcFailedCallback))

	// Success Test
	tcFile, tcCloser, err := test.OpenTCFile("TestSyncRead", tcFileSz)
	a.NoError(err)
	defer tcCloser()

	tcSuccessCallback := func(success bool) {
		a.True(success)
	}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		var testedBuffer [test.BufferSz]byte

		a.NoError(SyncRead(tcFile, tcOffset, testedBuffer[:], tcSuccessCallback))
		a.EqualValues(test.Buffer, testedBuffer[:])
	}
}

func TestAsyncRead(t *testing.T) {
	a := assert.New(t)

	// Fail Test
	tcFailedCallback := func(success bool) {
		a.False(success)
	}

	a.NoError(AsyncRead(nil, 0, nil, tcFailedCallback))

	// Success Test
	tcFile, tcCloser, err := test.OpenTCFile("TestAsyncRead", tcFileSz)
	a.NoError(err)
	defer tcCloser()

	testedReadOffset := int64(0)
	chBlock := make(chan bool)
	defer func() { close(chBlock) }()
	chDone := make(chan bool)
	defer func() { close(chDone) }()
	tcSuccessCallback := func(success bool) {
		a.True(success)

		// Wait until offset test result
		<-chBlock
		testedReadOffset += test.BufferSz
		chDone <- true
	}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		var testedBuffer [test.BufferSz]byte

		expectedOffset := tcOffset + test.BufferSz

		a.NoError(AsyncRead(tcFile, tcOffset, testedBuffer[:], tcSuccessCallback))
		a.NotEqual(expectedOffset, testedReadOffset)

		chBlock <- false
		// Wait until read done
		<-chDone

		a.EqualValues(test.Buffer, testedBuffer[:])
		a.Equal(expectedOffset, testedReadOffset)
	}
}
