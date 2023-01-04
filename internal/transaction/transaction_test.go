package transaction

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/pkg/measure"
	"math"
	"math/big"
	"os"
	"testing"
	"time"
)

func tcMakeIOList(tcFunc io.Type, jobId int64) []*io.IO {
	ios := make([]*io.IO, 0)

	for offset := int64(0); offset < (4096 << 3); offset += 4096 {
		ios = append(ios, io.NewIO(tcFunc, jobId, offset, make([]byte, 4096)))
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

	tested := Transaction{jobId: vRand.Int64()}
	testedCounter := 0

	tcSyncIO := func(_ *os.File, _ int64, _ []byte, cb func(bool)) error {
		if testSleep.Nanoseconds() > 0 {
			time.Sleep(testSleep)
		}
		testedCounter++
		cb(testErr == nil)
		return testErr
	}

	tcAsyncIO := func(_ *os.File, _ int64, _ []byte, cb func(bool)) error {
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
