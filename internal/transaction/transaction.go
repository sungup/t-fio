package transaction

import (
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"github.com/sungup/t-fio/pkg/measure"
	"os"
	"sync"
)

type Transaction struct {
	jobId int64
	ios   []*io.IO
	fp    *os.File
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
		if err = item.Issue(t.fp, wait); err != nil {
			break
		}
	}

	// wait until all transaction completed
	wait.Wait()

	return err
}

func (t *Transaction) AddIO(ioType io.Type, offset int64, buffer *bytebuf.ByteBuf) {
	t.ios = append(t.ios, io.New(ioType, t.jobId, offset, buffer))
}

func (t *Transaction) IOs() int {
	return len(t.ios)
}

func NewTransaction(jobId int64, fp *os.File) *Transaction {
	return &Transaction{jobId: jobId, fp: fp}
}
