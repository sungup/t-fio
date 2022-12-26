package rand

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZipf_Uint64(t *testing.T) {
	nRange := uint64(1000)
	buckets := make([]uint64, nRange)
	tested, _ := NewZipf(0, nRange, 0, 1.2)
	loop := 1000000

	tested.EnableHash(false)

	for i := 0; i < loop; i++ {
		testedV := tested.Uint64()
		assert.Less(t, testedV, nRange)
		buckets[testedV]++
	}

	last := buckets[0]
	for i, testedC := range buckets[1:] {
		// TODO is this the best way???
		// Top 3% values should be decrementing order because there are a lot the large appearance
		// change differentials.
		if i < 30 {
			assert.Less(t, testedC, last)
		} else {
			break
		}

		last = testedC
	}
}

func TestNewZipf(t *testing.T) {
	// initializing fail
	tested, err := NewZipf(0, 100, -2.0, 1.2)
	assert.Nil(t, tested)
	assert.EqualError(t, err, fmt.Sprintf("unexpected center range: %v", -2.0))

	// zipf creation fail
	tested, err = NewZipf(0, 100, 0.5, 0.5)
	assert.Nil(t, tested)
	assert.EqualError(t, err, fmt.Sprintf("theta value is not acceptable to create zipf random: %v", 0.5))

	// successful creation
	tested, err = NewZipf(0, 100, 0.5, 1.2)
	assert.NotNil(t, tested)
	assert.NoError(t, err)
}
