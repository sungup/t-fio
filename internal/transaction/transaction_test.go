package transaction

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/test"
	"math"
	"math/big"
	"testing"
	"time"
)

func tcMakeIOList(tcFunc engine.DoIO, jobId int64) (list []*io.IO, closer func()) {
	ios := make([]*io.IO, 0)

	for offset := int64(0); offset < (4096 << 3); offset += 4096 {
		ios = append(ios, io.New(tcFunc, jobId, offset, bytebuf.Alloc(4096)))
	}

	return ios, bytebuf.ForceCleanByteBufPool
}

func TestTransaction_ProcessAll(t *testing.T) {
	const errMessage = "TestTransaction_ProcessAll error message"
	var (
		testErr   error
		testSleep time.Duration
		closer    func()
	)

	vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	tested := Transaction{jobId: vRand.Int64()}
	testedCounter := 0

	tcSyncIO := func(_ []byte, _ int64, cb engine.Callback) error {
		if testSleep.Nanoseconds() > 0 {
			time.Sleep(testSleep)
		}
		testedCounter++
		cb(test.BufferSz, nil)
		return testErr
	}

	tcAsyncIO := func(_ []byte, _ int64, cb engine.Callback) error {
		go func() {
			if testSleep.Nanoseconds() > 0 {
				time.Sleep(testSleep)
			}
			testedCounter++
			cb(test.BufferSz, nil)
		}()
		return testErr
	}

	// Normal condition check for sync IO
	testSleep = time.Millisecond * 100
	testErr = nil
	testedCounter = 0
	lat := measure.LatencyMeasureStart()
	tested.ios, closer = tcMakeIOList(tcSyncIO, tested.jobId)
	assert.NoError(t, tested.ProcessAll())
	assert.Greater(t, lat().Nanoseconds(), testSleep.Nanoseconds()*int64(len(tested.ios)))
	assert.Equal(t, len(tested.ios), testedCounter)
	closer()

	// Normal condition check for async IO
	testSleep = time.Millisecond * 100
	testErr = nil
	testedCounter = 0
	lat = measure.LatencyMeasureStart()
	tested.ios, closer = tcMakeIOList(tcAsyncIO, tested.jobId)
	assert.NoError(t, tested.ProcessAll())
	assert.Greater(t, lat(), testSleep)
	assert.Less(t, lat(), testSleep*2)
	assert.Equal(t, len(tested.ios), testedCounter)
	closer()

	// Error condition check for sync IO
	testSleep = 0
	testErr = fmt.Errorf(errMessage)
	testedCounter = 0
	tested.ios, closer = tcMakeIOList(tcSyncIO, tested.jobId)
	assert.EqualError(t, tested.ProcessAll(), errMessage)
	assert.Equal(t, testedCounter, 1)
	closer()

	// Error condition check for async IO
	testSleep = 0
	testErr = fmt.Errorf(errMessage)
	testedCounter = 0
	tested.ios, closer = tcMakeIOList(tcAsyncIO, tested.jobId)
	assert.EqualError(t, tested.ProcessAll(), errMessage)
	assert.Equal(t, testedCounter, 1)
	closer()
}

func TestTransaction_AddIO(t *testing.T) {
	tested := &Transaction{
		jobId: 0,
		ios:   make([]*io.IO, 0),
	}

	tcFunc := func(_ []byte, _ int64, _ engine.Callback) (err error) { return nil }

	for sz := 1; sz <= 1024; sz++ {
		tested.AddIO(tcFunc, 0, bytebuf.Alloc(4096))
		assert.Len(t, tested.ios, sz)
	}

	bytebuf.ForceCleanByteBufPool()
}

func TestTransaction_IOs(t *testing.T) {
	tested := &Transaction{
		jobId: 0,
		ios:   make([]*io.IO, 0),
	}

	tcFunc := func(p []byte, offset int64, callback engine.Callback) (err error) { return nil }

	for expectedSz := 1; expectedSz <= 1024; expectedSz++ {
		tested.ios = append(tested.ios, io.New(tcFunc, tested.jobId, 0, bytebuf.Alloc(1024)))
		assert.Equal(t, expectedSz, tested.IOs())
	}

	bytebuf.ForceCleanByteBufPool()
}

func TestNewTransaction(t *testing.T) {
	const loop = 1000

	for jobId := 0; jobId < loop; jobId++ {
		generated := NewTransaction(int64(jobId))
		assert.Equal(t, int64(jobId), generated.jobId)
		assert.Empty(t, generated.ios)
	}
}
