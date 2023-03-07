package io

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/test"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"
)

func tcMakeIOStruct(issue engine.DoIO) *IO {
	tc := &IO{
		jobId:  0,
		offset: time.Now().UnixNano(),
		buffer: bytebuf.Alloc(test.BufferSz),
		issue:  issue,
	}

	return tc
}

func TestIO_Issue(t *testing.T) {
	var (
		tested        *IO
		expectedError error

		expectedWait = &sync.WaitGroup{}
	)

	testedIO := func(testedBuf []byte, testedOffset int64, testedCB engine.Callback) (err error) {
		assert.NotNil(t, testedBuf)
		assert.Equal(t, tested.buffer.Buffer(), testedBuf)
		assert.Equal(t, tested.offset, testedOffset)
		assert.NotNil(t, testedCB)

		testedCB(test.BufferSz, expectedError)

		return expectedError
	}
	defer func() { bytebuf.ForceCleanByteBufPool() }()

	// IssueDeprecated error test
	assert.NotPanics(t, func() {
		expectedError = fmt.Errorf("error data")
		tested = tcMakeIOStruct(testedIO)
		expectedWait.Add(1)
		assert.EqualError(t, tested.Issue(expectedWait), expectedError.Error())
	})

	// IssueDeprecated no-error test
	assert.NotPanics(t, func() {
		expectedError = nil
		tested = tcMakeIOStruct(testedIO)
		expectedWait.Add(1)
		assert.NoError(t, tested.Issue(expectedWait))
	})
}

func TestNew(t *testing.T) {
	tcTypes := []engine.Engine{&engine.AsyncIO{}, &engine.SyncIO{}}

	for _, tcEngine := range tcTypes {
		for _, tc := range []engine.DoIO{tcEngine.ReadAt, tcEngine.WriteAt} {
			vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			buffer := bytebuf.Alloc(4096)
			tested := New(tc, vRand.Int64(), vRand.Int64()+1, buffer)

			assert.Equal(t, vRand.Int64(), tested.jobId)
			assert.Equal(t, vRand.Int64()+1, tested.offset)
			assert.Equal(t, buffer, tested.buffer)
			assert.NotNil(t, tested.issue)
		}
	}
}
