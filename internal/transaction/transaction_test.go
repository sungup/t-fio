package transaction

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/test"
	"math"
	"math/big"
	"os"
	"testing"
	"time"
)

func tcMakeIOList(tcFunc io.Type, jobId int64) []*io.IO {
	ios := make([]*io.IO, 0)

	for offset := int64(0); offset < (4096 << 3); offset += 4096 {
		ios = append(ios, io.New(tcFunc, jobId, offset, make([]byte, 4096)))
	}

	return ios
}

func TestTransaction_ProcessAll(t *testing.T) {
	const errMessage = "TestTransaction_ProcessAll error message"
	var (
		testErr   error
		testSleep time.Duration
	)

	vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	testFP, testFPClose, _ := test.OpenTCFile("TestTransaction_ProcessAll", 4096)
	defer testFPClose()

	tested := Transaction{jobId: vRand.Int64(), fp: testFP}
	testedCounter := 0

	tcSyncIO := func(fp *os.File, _ int64, _ []byte, cb func(bool)) error {
		assert.Equal(t, testFP, fp)
		if testSleep.Nanoseconds() > 0 {
			time.Sleep(testSleep)
		}
		testedCounter++
		cb(testErr == nil)
		return testErr
	}

	tcAsyncIO := func(fp *os.File, _ int64, _ []byte, cb func(bool)) error {
		assert.Equal(t, testFP, fp)
		go func() {
			if testSleep.Nanoseconds() > 0 {
				time.Sleep(testSleep)
			}
			testedCounter++
			cb(testErr == nil)
		}()
		return testErr
	}

	// Normal condition check for sync IO
	testSleep = time.Millisecond * 100
	testErr = nil
	testedCounter = 0
	lat := measure.LatencyMeasureStart()
	tested.ios = tcMakeIOList(tcSyncIO, tested.jobId)
	assert.NoError(t, tested.ProcessAll())
	assert.Greater(t, lat().Nanoseconds(), testSleep.Nanoseconds()*int64(len(tested.ios)))
	assert.Equal(t, len(tested.ios), testedCounter)

	// Normal condition check for async IO
	testSleep = time.Millisecond * 100
	testErr = nil
	testedCounter = 0
	lat = measure.LatencyMeasureStart()
	tested.ios = tcMakeIOList(tcAsyncIO, tested.jobId)
	assert.NoError(t, tested.ProcessAll())
	assert.Greater(t, lat(), testSleep)
	assert.Less(t, lat(), testSleep*2)
	assert.Equal(t, len(tested.ios), testedCounter)

	// Error condition check for sync IO
	testSleep = 0
	testErr = fmt.Errorf(errMessage)
	testedCounter = 0
	tested.ios = tcMakeIOList(tcSyncIO, tested.jobId)
	assert.EqualError(t, tested.ProcessAll(), errMessage)
	assert.Equal(t, testedCounter, 1)

	// Error condition check for async IO
	testSleep = 0
	testErr = fmt.Errorf(errMessage)
	testedCounter = 0
	tested.ios = tcMakeIOList(tcAsyncIO, tested.jobId)
	assert.EqualError(t, tested.ProcessAll(), errMessage)
	assert.Equal(t, testedCounter, 1)
}

func TestTransaction_AddIO(t *testing.T) {
	tested := &Transaction{
		jobId: 0,
		ios:   make([]*io.IO, 0),
		fp:    nil,
	}

	tcFunc := func(_ *os.File, _ int64, _ []byte, _ func(bool)) error { return nil }

	for sz := 1; sz <= 1024; sz++ {
		tested.AddIO(tcFunc, 0, make([]byte, 4096))
		assert.Len(t, tested.ios, sz)
	}
}

func TestNewTransaction(t *testing.T) {
	const loop = 1000
	closers := make([]func(), loop)
	tcFP := make([]*os.File, loop)

	for i := 0; i < loop; i++ {
		tcFP[i], closers[i], _ = test.OpenTCFile(fmt.Sprintf("TestNewTransaction-%d", i), 4096)
	}
	defer func() {
		for _, closer := range closers {
			closer()
		}
	}()

	for jobId, fp := range tcFP {
		generated := NewTransaction(int64(jobId), fp)

		assert.Equal(t, int64(jobId), generated.jobId)
		assert.Equal(t, fp, generated.fp)
	}
}
