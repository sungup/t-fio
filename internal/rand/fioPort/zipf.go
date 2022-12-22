package fioPort

import (
	"math"
	"math/rand"
)

const (
	zipfMaxGen = float64(10000000)
)

func hash(v uint64) uint64 {
	return v
}

// FIOSharedZipf is a shared structure between zipf and pareto rand
type FIOSharedZipf struct {
	nRange      uint64
	rnd         *rand.Rand
	randOff     uint64
	disableHash bool
}

func (z *FIOSharedZipf) sharedRandInit(nRange uint64, center float64, seed int64) {
	z.nRange = nRange
	z.rnd = rand.New(rand.NewSource(seed))

	z.randOff = rand.Uint64()

	if center != -1 {
		z.randOff = uint64(float64(z.nRange) * center)
	}
}

func (z *FIOSharedZipf) Uint64(v uint64) uint64 {
	if !z.disableHash {
		v = hash(v)
	}

	return (v + z.randOff) % z.nRange
}

// FIOZipf is the implemented Zipf randomizer in FIO
type FIOZipf struct {
	// common factor
	FIOSharedZipf

	// zipf factor
	theta float64
	zeta2 float64
	zetaN float64
}

func (z *FIOZipf) update() {
	toGen := int(math.Min(float64(z.nRange), zipfMaxGen))

	/*
	 * It can become very costly to generate long sequences. Just cap it at
	 * 10M max, that should be double in 1-2s on even slow machines.
	 * Precision will take a slight hit, but nothing major.
	 */
	for i := 0; i < toGen; i++ {
		z.zetaN += math.Pow(1.0/float64(i+1), z.theta)
	}
}

func (z *FIOZipf) init(nRange uint64, theta, center float64, seed int64) {
	z.sharedRandInit(nRange, center, seed)

	z.theta = theta
	z.zeta2 = math.Pow(1.0, z.theta) + math.Pow(0.5, z.theta)

	z.update()
}

func (z *FIOZipf) Uint64() uint64 {
	var (
		alpha, eta, randUni, randZ float64

		n   = z.nRange
		val uint64
	)

	alpha = 1.0 / (1.0 - z.theta)
	eta = (1.0 - math.Pow(2.0/float64(n), 1.0-z.theta)) / (1.0 - z.zeta2/z.zetaN)
	randUni = z.rnd.Float64() / math.MaxUint32
	randZ = randUni * z.zetaN

	switch {
	case randZ < 1.0:
		val = 1
	case randZ < (1.0 + math.Pow(0.5, z.theta)):
		val = 2
	default:
		val = 1 + uint64(float64(n)*math.Pow(eta*randUni-eta+1.0, alpha))
	}

	return z.FIOSharedZipf.Uint64(val - 1)
}

// FIOPareto is the implemented Pareto randomizer in FIO
type FIOPareto struct {
	// common factor
	FIOSharedZipf

	// ParetoFactor
	paretoPow float64
}

func (p *FIOPareto) init(nRange uint64, h, center float64, seed int64) {
	p.sharedRandInit(nRange, center, seed)
	p.paretoPow = math.Log(h) / math.Log(1.0-h)
}

func (p *FIOPareto) Uint64() uint64 {
	rnd := p.rnd.Float64() / math.MaxUint32

	return p.FIOSharedZipf.Uint64(uint64(float64(p.nRange-1) * math.Pow(rnd, p.paretoPow)))
}
