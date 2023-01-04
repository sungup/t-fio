package worker

import (
	"context"
	"github.com/sungup/t-fio/internal/transaction"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Worker struct {
	queue <-chan *transaction.Transaction
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case tr := <-w.queue:
			_ = tr.ProcessAll()

		case <-ctx.Done():
			return
		}
	}
}
