package rand

import (
	"fmt"
	"math/rand"
)

type Zipf struct {
	randomizer
	zipf *rand.Zipf
}

func (z *Zipf) Uint64() uint64 {
	return z.randomizer.hash(z.zipf.Uint64())
}

func NewZipf(seed int64, nRange uint64, center, theta float64) (Randomizer, error) {
	z := &Zipf{}
	if err := z.init(seed, nRange, center); err != nil {
		return nil, err
	}

	if z.zipf = rand.NewZipf(z.rand, theta, 1.0, nRange); z.zipf == nil {
		return nil, fmt.Errorf("theta value is not acceptable to create zipf random: %v", theta)
	}

	return z, nil
}
