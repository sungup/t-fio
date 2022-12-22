package io

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"os"
	"testing"
	"time"
)

func tcMakeIOStruct(ch chan<- *IO, issue func(*os.File, int64, []byte, func(bool)) error) *IO {
	tc := &IO{
		id:      time.Now().UnixNano(),
		offset:  time.Now().UnixNano(),
		buffer:  make([]byte, test.BufferSz),
		next:    nil,
		success: false,
		issue:   issue,
		ch:      ch,
	}

	tc.next = tc // self link to error test

	return tc
}

func TestIO_Issue(t *testing.T) {
	a := assert.New(t)

	var (
		tested        *IO
		expectedError          = fmt.Errorf("error data")
		expectedFP    *os.File = nil
	)

	tested = tcMakeIOStruct(nil, func(testedFP *os.File, testedOffset int64, testedBuf []byte, testedCB func(success bool)) error {
		a.Equal(expectedFP, testedFP)
		a.Equal(tested.offset, testedOffset)
		a.Equal(tested.buffer, testedBuf)

		return expectedError
	})

	// Issue error test
	next, err := tested.Issue(nil)
	a.EqualError(err, expectedError.Error())
	a.Nil(next)

	// Issue no-error test
	expectedError = nil
	next, err = tested.Issue(nil)
	a.NoError(err)
	a.Equal(next, tested.next)
}

func TestIO_Callback(t *testing.T) {
	a := assert.New(t)
	var (
		tcChan = make(chan *IO)
		tcIO   = tcMakeIOStruct(tcChan, nil)
		tested *IO
	)

	// True Test
	for _, expected := range []bool{true, false} {
		go tcIO.Callback(expected)
		tested = <-tcChan
		a.Equal(tcIO, tested)
		a.Equal(expected, tested.success)
	}
}
