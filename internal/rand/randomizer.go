package rand

import (
	"fmt"
	"github.com/sungup/t-fio/internal/hash"
	"math/rand"
)

type Randomizer interface {
	Uint64() uint64
	EnableHash(enable bool)
}

type randomizer struct {
	nRange      uint64
	rand        *rand.Rand
	randOff     uint64
	disableHash bool
}

func (r *randomizer) init(seed int64, nRange uint64, center float64) error {
	if center != -1 && (center < 0 || 1 < center) {
		return fmt.Errorf("unexpected center range: %v", center)
	}

	r.nRange = nRange
	r.rand = rand.New(rand.NewSource(seed))
	r.randOff = r.rand.Uint64() % nRange
	r.disableHash = false

	if center != -1 {
		r.randOff = uint64(float64(r.nRange) * center)
	}

	return nil
}

func (r *randomizer) EnableHash(enable bool) {
	r.disableHash = !enable
}

func (r *randomizer) hash(v uint64) uint64 {
	if !r.disableHash {
		v = hash.Hash(v)
	}

	return (r.randOff + v) % r.nRange
}
