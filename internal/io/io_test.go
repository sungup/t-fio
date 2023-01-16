package io

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/test"
	"math"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"
)

func tcMakeIOStruct(issue func(*os.File, int64, []byte, func(bool)) error) *IO {
	tc := &IO{
		jobId:  0,
		offset: time.Now().UnixNano(),
		buffer: bytebuf.Alloc(test.BufferSz),
		issue:  issue,
	}

	return tc
}

func TestIO_Issue(t *testing.T) {
	a := assert.New(t)

	var (
		tested       *IO
		expectedFP   *os.File = nil
		expectedWait          = &sync.WaitGroup{}

		expectedError = fmt.Errorf("error data")
	)

	tested = tcMakeIOStruct(func(testedFP *os.File, testedOffset int64, testedBuf []byte, testedCB func(success bool)) error {
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

		tested = tcMakeIOStruct(nil)

		tcIOLat = make([]func() time.Duration, 10)
	)
	defer func() { bytebuf.ForceCleanByteBufPool() }()

	// True Test
	totalLat := measure.LatencyMeasureStart()
	for i := range tcIOLat {
		jobWait.Add(1)
		lat := measure.LatencyMeasureStart()
		go func(lat func() time.Duration) {
			tested = tcMakeIOStruct(nil)
			tested.wait = jobWait
			tested.latency = lat
			time.Sleep(tcSleep)
			tested.Callback(true)
		}(lat)
		tcIOLat[i] = lat
	}

	// check all job has been done after tcSleep
	go func() {
		tcWait.Add(1)
		jobWait.Wait()
		assert.Greater(t, totalLat(), tcSleep)
		tcWait.Done()
	}()

	// check each job has been done after tcSleep
	for _, lat := range tcIOLat {
		go func(lat func() time.Duration) {
			tcWait.Add(1)
			jobWait.Wait()
			assert.Greater(t, lat(), tcSleep)
			tcWait.Done()
		}(lat)
	}

	// all check thread has been launched in 500msec
	assert.Less(t, totalLat(), tcSleep)

	tcWait.Wait()
}

func TestNew(t *testing.T) {
	tcTypes := []Type{SyncRead, AsyncRead, Write}

	for _, tc := range tcTypes {
		vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		buffer := bytebuf.Alloc(4096)
		tested := New(tc, vRand.Int64(), vRand.Int64()+1, buffer)

		assert.Equal(t, vRand.Int64(), tested.jobId)
		assert.Equal(t, vRand.Int64()+1, tested.offset)
		assert.Equal(t, buffer, tested.buffer)
		assert.NotNil(t, tested.issue)
	}
}
