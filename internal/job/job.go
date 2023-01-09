package job

import (
	"context"
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/internal/pattern"
	"github.com/sungup/t-fio/internal/transaction"
	"os"
	"time"
)

type Job struct {
	fp       *os.File
	jobId    int64
	ioType   io.Type
	ioSize   int
	address  *pattern.Generator
	delay    time.Duration
	trLength int
	buffer   chan *transaction.Transaction

	newBuffer func(size int) []byte
}

func (j *Job) newTransaction() *transaction.Transaction {
	tr := transaction.NewTransaction(j.jobId, j.fp)

	for i := 0; i < j.trLength; i++ {
		tr.AddIO(j.ioType, j.address.Offset(), j.newBuffer(j.ioSize))
	}

	return tr
}

func (j *Job) Run(ctx context.Context) {
	// wait until task is ready
	if j.delay > 0 {
		time.Sleep(j.delay)
	}

	for {
		select {
		case j.buffer <- j.newTransaction():
		case <-ctx.Done():
			return
		}
	}
}

func (j *Job) TransactionReceiver() <-chan *transaction.Transaction {
	return j.buffer
}
