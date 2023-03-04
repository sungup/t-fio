package job

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/pattern"
	"github.com/sungup/t-fio/internal/transaction"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/pkg/sys"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func tcMakePatternGenerator() (g *pattern.Generator) {
	jsonStr := "{}"
	opts := pattern.Options{}
	_ = json.Unmarshal([]byte(jsonStr), &opts)

	g, _ = opts.MakeGenerator()

	return g
}

func TestJob_newTransaction(t *testing.T) {
	const loop = 1000
	var (
		expectedJobId  = rand.Int63()
		expectedIoType = func(_ sys.File, _ int64, _ []byte, _ func(bool)) error { return nil }
		expectedTRLen  = 16
	)

	tested := Job{
		fp:        nil,
		jobId:     expectedJobId,
		ioType:    expectedIoType,
		ioSize:    4096,
		address:   tcMakePatternGenerator(),
		delay:     0,
		trLength:  expectedTRLen,
		buffer:    make(chan *transaction.Transaction, 1),
		newBuffer: AllocReadBuffer,
	}

	for i := 0; i < loop; i++ {
		generated := tested.newTransaction()
		assert.NotNil(t, generated)
		assert.Equal(t, expectedTRLen, generated.IOs())
	}
}

func TestJob_Run(t *testing.T) {
	const loop = 1000

	tested := Job{
		fp:        nil,
		jobId:     rand.Int63(),
		ioType:    func(_ sys.File, _ int64, _ []byte, _ func(bool)) error { return nil },
		ioSize:    4096,
		address:   tcMakePatternGenerator(),
		delay:     0,
		trLength:  1,
		buffer:    nil,
		newBuffer: AllocReadBuffer,
	}

	tcDelay := time.Millisecond * 100
	tcWaitCancel := time.Millisecond * 10
	tcDeadline := tcDelay + time.Millisecond*200
	wg := sync.WaitGroup{}

	// 1. Delay Test
	tested.delay = tcDelay
	tested.buffer = make(chan *transaction.Transaction, 1)
	wg.Add(1)
	start := time.Now()
	timeout, cancel := context.WithTimeout(context.Background(), tcDeadline)
	go func() {
		tested.Run(timeout)
		wg.Done()
	}()
	cancel()
	wg.Wait()
	assert.WithinRange(t, time.Now(), start.Add(tcDelay), start.Add(tcDelay+tcWaitCancel))

	// 1. Normal Test
	tested.delay = 0
	tested.buffer = make(chan *transaction.Transaction, 1)
	wg.Add(1)
	start = time.Now()
	timeout, cancel = context.WithTimeout(context.Background(), time.Second)
	go func() {
		tested.Run(timeout)
		wg.Done()
	}()

	// all transaction should be done until tcDeadline (300msec)
	for i := 0; i < loop; i++ {
		assert.IsType(t, &transaction.Transaction{}, <-tested.buffer)
		assert.WithinRange(t, time.Now(), start, start.Add(tcDeadline))
	}
	time.Sleep(tcWaitCancel)
	assert.NotEmpty(t, tested.buffer)
	cancel()
	wg.Wait()
	assert.WithinRange(t, time.Now(), start, start.Add(tcDeadline))
}

func TestJob_TransactionReceiver(t *testing.T) {
	const loop = 1000
	const deadline = time.Second
	tcQueue := make(chan *transaction.Transaction, 1)
	tested := Job{
		fp:        nil,
		jobId:     0,
		ioType:    func(_ sys.File, _ int64, _ []byte, _ func(bool)) error { return nil },
		ioSize:    1024,
		address:   tcMakePatternGenerator(),
		delay:     0,
		trLength:  4,
		buffer:    tcQueue,
		newBuffer: AllocReadBuffer,
	}

	generated := tested.TransactionReceiver()
	assert.NotNil(t, generated)

	// Run job creation routine in background to check returned channel is working
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	go tested.Run(ctx)

	// Run job receiving test loop times
	count := 0
	wait := sync.WaitGroup{}
	wait.Add(1)

	lat := measure.LatencyMeasureStart()
	go func(wait *sync.WaitGroup) {
		defer wait.Done()

		c, closer := context.WithTimeout(context.Background(), deadline)
		defer closer()

		for {
			select {
			case tr := <-generated:
				assert.IsType(t, &transaction.Transaction{}, tr)
				count++
				if count == loop {
					return
				}
			case <-c.Done():
				return
			}
		}

	}(&wait)

	wait.Wait()

	assert.Equal(t, loop, count)
	assert.Less(t, lat(), deadline)
}
