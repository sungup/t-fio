package transaction

import (
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"sync"
)

type Transaction struct {
	jobId int64
	ios   []*io.IO
}

func (t *Transaction) ProcessAll() (err error) {
	lat := measure.LatencyMeasureStart()
	defer func() {
		// TODO call stat collector function
		lat()
	}()

	wait := &sync.WaitGroup{}

	for _, item := range t.ios {
		wait.Add(1)
		if err = item.Issue(wait); err != nil {
			break
		}
	}

	// wait until all transaction completed
	wait.Wait()

	return err
}

func (t *Transaction) AddIO(ioAction engine.DoIO, offset int64, buffer *bytebuf.ByteBuf) {
	t.ios = append(t.ios, io.New(ioAction, t.jobId, offset, buffer))
}

func (t *Transaction) IOs() int {
	return len(t.ios)
}

func NewTransaction(jobId int64) *Transaction {
	return &Transaction{jobId: jobId}
}
