package job

import (
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/internal/pattern"
	"github.com/sungup/t-fio/internal/transaction"
	"os"
	"time"
)

type Job struct {
	fp      *os.File
	jobId   int64
	address *pattern.Generator
	ioType  io.Type
	ioSize  int
	delay   time.Duration

	trLength int
}

func (j *Job) MakeTransaction() *transaction.Transaction {
	tr := transaction.NewTransaction(j.jobId, j.fp)

	for i := 0; i < j.trLength; i++ {
		// TODO
	}

	return tr
}
