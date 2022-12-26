package pattern

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequencer_PageNo(t *testing.T) {
	seqTotal := int64(1024)
	loop := seqTotal * seqTotal

	tested := Sequencer{
		until:  seqTotal,
		cursor: 0,
	}

	for i := int64(0); i < loop; i++ {
		assert.Equal(t, i%seqTotal, tested.PageNo())
	}
}

func TestNewSequencer(t *testing.T) {
	loop := int64(1024)

	for expected := int64(0); expected < loop; expected++ {
		tested := NewSequencer(expected)
		assert.Zero(t, tested.(*Sequencer).cursor)
		assert.Equal(t, expected, tested.(*Sequencer).until)
	}
}
