package transaction

import (
	"github.com/sungup/t-fio/internal/io"
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
			return err
		}
	}

	// wait until all transaction completed
	wait.Wait()

	return nil
}
