package rand

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/hash"
	"math/rand"
	"testing"
	"time"
)

func TestRandomizer_init(t *testing.T) {
	a := assert.New(t)

	tested := randomizer{}

	// error raise check
	for _, center := range []float64{-2.0, -1.1, -0.5, 1.001} {
		a.Error(tested.init(0, 0, center), fmt.Errorf("unexpected center range: %v", center))
	}

	rangeLoop := uint64(1024 * 1024)
	seedLoop := int64(10)
	for expectedRange := uint64(1); expectedRange < rangeLoop; expectedRange <<= 1 {
		for expectedSeed := int64(0); expectedSeed < seedLoop; expectedSeed++ {
			offset := rand.New(rand.NewSource(expectedSeed)).Uint64() % expectedRange

			a.NoError(tested.init(expectedSeed, expectedRange, -1))
			a.Equal(expectedRange, tested.nRange)
			a.NotNil(tested.rand)
			a.Equal(offset, tested.randOff)
			a.False(tested.disableHash)

			for _, center := range []float64{0.0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0} {
				expectedOffset := uint64(float64(expectedRange) * center)

				a.NoError(tested.init(expectedSeed, expectedRange, center))
				a.Equal(expectedRange, tested.nRange)
				a.NotNil(tested.rand)
				a.Equal(expectedOffset, tested.randOff)
				a.False(tested.disableHash)
			}
		}
	}
}

func TestRandomizer_EnableHash(t *testing.T) {
	tested := randomizer{disableHash: false}

	for _, in := range []bool{true, false} {
		expected := !in

		tested.EnableHash(in)
		assert.Equal(t, expected, tested.disableHash)
	}
}

func TestRandomizer_hash(t *testing.T) {
	a := assert.New(t)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	nRange := uint64(1024 * 128)
	center := 0.5
	offset := uint64(float64(nRange) * center)

	tested := randomizer{}
	_ = tested.init(0, nRange, center)

	for i := uint64(0); i < nRange; i++ {
		in := rnd.Uint64() % nRange

		tested.EnableHash(true)
		expected := (offset + hash.Hash(in)) % nRange
		a.Less(tested.hash(in), nRange)
		a.Equal(expected, tested.hash(in))

		tested.EnableHash(false)
		expected = (offset + in) % nRange
		a.Less(tested.hash(in), nRange)
		a.Equal(expected, tested.hash(in))
	}

}
