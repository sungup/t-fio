package rand

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequencer_PageNo(t *testing.T) {
	base := int64(1024)
	seqSz := int64(1024)
	loop := seqSz * seqSz

	tested := Sequencer{
		base:   base,
		until:  base + seqSz,
		cursor: base,
	}

	for i := int64(0); i < loop; i++ {
		assert.Equal(t, i%seqSz+base, tested.PageNo())
	}
}

func TestNewSequencer(t *testing.T) {
	pageSz := int64(512)           // base sector size
	baseOffset := pageSz * 8       // 4KB page offset
	ioRangeSz := pageSz * 8 * 4096 // 16MB IO Range

	expectedBase := baseOffset / pageSz
	expectedUntil := (baseOffset + ioRangeSz) / pageSz
	expectedCursor := expectedBase

	tested := NewSequencer(pageSz, baseOffset, ioRangeSz)
	assert.Equal(t, expectedBase, tested.base)
	assert.Equal(t, expectedUntil, tested.until)
	assert.Equal(t, expectedCursor, tested.cursor)
}
