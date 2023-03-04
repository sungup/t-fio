package io

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/pkg/sys"
	"github.com/sungup/t-fio/test"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"
)

func tcMakeIOStruct(issue DoIO) *IO {
	tc := &IO{
		jobId:  0,
		offset: time.Now().UnixNano(),
		buffer: bytebuf.Alloc(test.BufferSz),
		issue:  issue,
	}

	return tc
}

func tcMakeIOStructDeprecated(issue func(sys.File, int64, []byte, func(bool)) error) *IO {
	tc := &IO{
		jobId:           0,
		offset:          time.Now().UnixNano(),
		buffer:          bytebuf.Alloc(test.BufferSz),
		issueDeprecated: issue,
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

		return expectedError
	}
	defer func() { bytebuf.ForceCleanByteBufPool() }()

	// Issue error test
	expectedError = fmt.Errorf("error data")
	tested = tcMakeIOStruct(testedIO)
	assert.EqualError(t, tested.Issue2(expectedWait), expectedError.Error())

	// Issue no-error test
	expectedError = nil
	tested = tcMakeIOStruct(testedIO)
	assert.NoError(t, tested.Issue2(expectedWait))
}

func TestIO_IssueDeprecated(t *testing.T) {
	a := assert.New(t)

	var (
		tested       *IO
		expectedFP   sys.File = nil
		expectedWait          = &sync.WaitGroup{}

		expectedError = fmt.Errorf("error data")
	)

	tested = tcMakeIOStructDeprecated(func(testedFP sys.File, testedOffset int64, testedBuf []byte, testedCB func(success bool)) error {
		a.Equal(expectedFP, testedFP)
		a.Equal(tested.offset, testedOffset)
		a.Equal(tested.buffer.Buffer(), testedBuf)

		a.Equal(expectedWait, tested.wait)
		a.NotNil(tested.latency)

		return expectedError
	})
	defer func() { bytebuf.ForceCleanByteBufPool() }()

	// Issue error test
	a.EqualError(tested.Issue(nil, expectedWait), expectedError.Error())

	// Issue no-error test
	expectedError = nil
	tested.wait = nil
	tested.latency = nil
	a.NoError(tested.Issue(nil, expectedWait))
}

func TestIO_Callback(t *testing.T) {
	var (
		tcWait  = &sync.WaitGroup{}
		jobWait = &sync.WaitGroup{}
		tcSleep = time.Millisecond * 500

		tested = tcMakeIOStructDeprecated(nil)

		tcIOLat = make([]func() time.Duration, 10)
	)
	defer func() { bytebuf.ForceCleanByteBufPool() }()

	// True Test
	totalLat := measure.LatencyMeasureStart()
	for i := range tcIOLat {
		jobWait.Add(1)
		lat := measure.LatencyMeasureStart()
		go func(lat func() time.Duration) {
			tested = tcMakeIOStructDeprecated(nil)
			tested.wait = jobWait
			tested.latency = lat
			time.Sleep(tcSleep)
			tested.Callback(true)
		}(lat)
		tcIOLat[i] = lat
	}

	// check all job has been done after tcSleep
	tcWait.Add(1)
	go func() {
		jobWait.Wait()
		assert.Greater(t, totalLat(), tcSleep)
		tcWait.Done()
	}()

	// check each job has been done after tcSleep
	for _, lat := range tcIOLat {
		tcWait.Add(1)
		go func(lat func() time.Duration) {
			jobWait.Wait()
			assert.Greater(t, lat(), tcSleep)
			tcWait.Done()
		}(lat)
	}

	// all check thread has been launched in 500msec
	assert.Less(t, totalLat(), tcSleep)

	tcWait.Wait()
	assert.Greater(t, totalLat(), tcSleep)
}

func TestNew(t *testing.T) {
	for _, tc := range []Type{SyncRead, AsyncRead, Write} {
		vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		buffer := bytebuf.Alloc(4096)
		tested := New(tc, vRand.Int64(), vRand.Int64()+1, buffer)

		assert.Equal(t, vRand.Int64(), tested.jobId)
		assert.Equal(t, vRand.Int64()+1, tested.offset)
		assert.Equal(t, buffer, tested.buffer)
		assert.NotNil(t, tested.issueDeprecated)
	}

	// Version2 TC
	tcTypes := []engine.Engine{&engine.AsyncIO{}, &engine.SyncIO{}}

	for _, tcEngine := range tcTypes {
		for _, tc := range []DoIO{tcEngine.ReadAt, tcEngine.WriteAt} {
			vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			buffer := bytebuf.Alloc(4096)
			tested := New2(tc, vRand.Int64(), vRand.Int64()+1, buffer)

			assert.Equal(t, vRand.Int64(), tested.jobId)
			assert.Equal(t, vRand.Int64()+1, tested.offset)
			assert.Equal(t, buffer, tested.buffer)
			assert.NotNil(t, tested.issue)
		}
	}
}
