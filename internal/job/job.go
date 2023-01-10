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
	fp       *os.File           // applied by constructor
	jobId    int64              // automatically assigned
	ioType   io.Type            // receive from Options
	ioSize   int                // receive from Options
	address  *pattern.Generator // created by pattern.Options
	delay    time.Duration      // receive from Options
	trLength int                // receive from Options

	buffer chan *transaction.Transaction // created at construction

	newBuffer func(size int) []byte // selected by ioType
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
