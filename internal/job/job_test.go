package job

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/internal/pattern"
	"github.com/sungup/t-fio/internal/transaction"
	"github.com/sungup/t-fio/pkg/measure"
	"math/rand"
	"os"
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
		expectedIoType = func(_ *os.File, _ int64, _ []byte, _ func(bool)) error { return nil }
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
		queue:     make(chan *transaction.Transaction, 1),
		newBuffer: io.AllocReadBuffer,
	}

	for i := 0; i < loop; i++ {
		generated := tested.newTransaction()
		assert.NotNil(t, generated)
		assert.Equal(t, expectedTRLen, generated.IOs())
	}
}

func TestJob_Run(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestJob_TransactionReceiver(t *testing.T) {
	tcQueue := make(chan *transaction.Transaction, 1)
	tested := Job{
		fp:        nil,
		jobId:     0,
		ioType:    func(_ *os.File, _ int64, _ []byte, _ func(bool)) error { return nil },
		ioSize:    1024,
		address:   nil,
		delay:     0,
		trLength:  4,
		queue:     tcQueue,
		newBuffer: nil,
	}

	generated := tested.TransactionReceiver()
	assert.NotNil(t, generated)

	count := 0
	wait := sync.WaitGroup{}
	wait.Add(1)

	lat := measure.LatencyMeasureStart()
	go func(wait *sync.WaitGroup) {
		for {
			select {
			case <-generated:
				count++
			case <-time.After(time.Millisecond * 100):
				wait.Done()
				break
			}
		}

	}(&wait)

	wait.Wait()

	t.Log(lat())
	assert.Fail(t, "not yet implemented")
}
