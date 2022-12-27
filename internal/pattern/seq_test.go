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

func TestSeqOptions_MakeIOPattern(t *testing.T) {
	tested := &SeqOptions{StartFrom: 0.5}

	loop := int64(16384)
	for expectedRange := int64(16); expectedRange < loop; expectedRange <<= 1 {
		expectedCursor := expectedRange >> 1

		testedV, err := tested.MakeIOPattern(expectedRange)
		assert.NoError(t, err)
		assert.Equal(t, expectedRange, testedV.(*Sequencer).until)
		assert.Equal(t, expectedCursor, testedV.(*Sequencer).cursor)
	}
}
