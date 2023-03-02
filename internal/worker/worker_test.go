package worker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/transaction"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"github.com/sungup/t-fio/pkg/sys"
	"github.com/sungup/t-fio/test"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorker_Run(t *testing.T) {
	// test scenario definition
	// 1. TestWorker_Run launch core / 2 count workers
	// 2. all workers should execute 500 transactions with async manner
	// 3. each transaction contains 32 unit ios
	const (
		ios   = 32
		loop  = 500
		sleep = time.Millisecond * 10

		expectedIssued = int32(loop * ios)
	)

	defer bytebuf.ForceCleanByteBufPool()

	workers := runtime.NumCPU() >> 1 // use only half of system cpu
	if workers == 0 {
		workers = 1
	}

	expectedMinRunTime := sleep * loop / time.Duration(workers)
	expectedMaxRunTime := time.Duration(float64(expectedMinRunTime.Nanoseconds()) * 1.1)

	queue := make(chan *transaction.Transaction)
	defer func() { close(queue) }()
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	testedWorkers := make([]*Worker, workers)
	for i := range testedWorkers {
		testedWorkers[i] = &Worker{queue: queue}
	}
	testedDone := 0
	testedIssued := int32(0)

	// make transaction data
	tcIOFunc := func(_ sys.File, _ int64, _ []byte, cb func(bool)) error {
		go func() {
			time.Sleep(sleep)
			atomic.AddInt32(&testedIssued, 1)
			cb(true)
		}()

		return nil
	}
	tcFP, closer, _ := test.OpenTCFile("TestWorker_Run", 10<<20)
	defer closer()
	tcTransactions := make([]*transaction.Transaction, loop)
	for jobId := range tcTransactions {
		tcTransactions[jobId] = transaction.NewTransaction(int64(jobId), tcFP)
		for i := 0; i < ios; i++ {
			tcTransactions[jobId].AddIO(tcIOFunc, int64(i*14), bytebuf.Alloc(4096))
		}
	}

	// executed multiple worker
	for _, tested := range testedWorkers {
		wg.Add(1)
		go func(tested *Worker) {
			tested.Run(ctx)
			testedDone++
			wg.Done()
		}(tested)
	}

	// send all transactions to workers
	testedRunTime := measure.LatencyMeasureStart()
	for _, tr := range tcTransactions {
		queue <- tr
	}
	assert.Less(t, testedRunTime(), expectedMaxRunTime)
	assert.Greater(t, testedRunTime(), expectedMinRunTime)
	assert.Zero(t, testedDone)

	// close context and wait until all worker routine are closed
	cancel()
	wg.Wait()

	assert.Equal(t, len(testedWorkers), testedDone) // all context closed safely
	assert.Equal(t, expectedIssued, testedIssued)   // all issued io has been completed
}
