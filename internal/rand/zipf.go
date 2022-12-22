package rand

import (
	"math/rand"
	"time"
)

type ZipfRand struct {
	base int64

	rand *rand.Zipf
}

func (z *ZipfRand) PageNo() int64 {
	return z.base + int64(z.rand.Uint64())
}

func NewZipf(pageSz, baseOffset, ioRangeSz int64, s, v float64) Randomizer {
	// FIO uses theta and value population
	// we can change value population to ranks from "ioRangeSz / pageSz * vPopulation", and also
	// we can change theta to

	return &ZipfRand{
		base: baseOffset / pageSz,
		rand: rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), s, v, uint64(ioRangeSz/pageSz)),
	}
}
