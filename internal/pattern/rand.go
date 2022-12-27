package pattern

import (
	"github.com/sungup/t-fio/internal/rand"
)

type Randomizer struct {
	rnd rand.Rand
}

func (r *Randomizer) PageNo() int64 {
	// simply casting type because lba number cannot exceed MaxInt63 value
	return int64(r.rnd.Uint64())
}

type RandOptions struct {
	rand.Options
}

func (r *RandOptions) MakeIOPattern(nRange int64) (pattern IOPattern, err error) {
	var rnd rand.Rand

	if rnd, err = r.MakeRandomizer(uint64(nRange)); err == nil {
		pattern = &Randomizer{rnd: rnd}
	}

	return pattern, err
}
