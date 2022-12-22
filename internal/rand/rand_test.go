package rand

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type tcSequencer struct {
	i int64
}

func (t *tcSequencer) PageNo() int64 {
	pageNo := t.i
	t.i++
	return pageNo
}

func newTCSequencer() *tcSequencer {
	return &tcSequencer{i: 0}
}

func TestRand_Offset(t *testing.T) {
	for pg := int64(512); pg < 4096; pg *= 2 {
		tested := Rand{rnd: newTCSequencer(), pageSz: pg}

		for i := int64(0); i < 4096; i++ {
			assert.Equal(t, pg*i, tested.Offset())
		}
	}
}
