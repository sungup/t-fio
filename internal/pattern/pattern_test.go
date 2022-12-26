package pattern

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

func TestRand_Offset(t *testing.T) {
	for pageSz := int64(512); pageSz < 4096; pageSz *= 2 {
		offset := pageSz * 1024
		ioRange := pageSz * 1024 * 1024
		end := offset + ioRange

		tested := Pattern{
			rnd:        &tcSequencer{i: 0},
			pageSz:     pageSz,
			pageOffset: offset,
		}

		for expected := offset; expected < end; expected += pageSz {
			assert.Equal(t, expected, tested.Offset())
		}
	}
}
