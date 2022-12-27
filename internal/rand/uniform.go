package rand

// Uniform is a random value generator with uniform distribution manner
type Uniform struct {
	core
}

// Uint64 is a random generating function but, there is no need hashing because all generated
// random values has same probability.
func (u *Uniform) Uint64() uint64 {
	return u.core.rand.Uint64() % u.nRange
}

// UniformOptions is an additional option container to generate to Uniform randomizer
type UniformOptions struct{}

func (o *UniformOptions) MakeRandomizer(seed int64, nRange uint64, center float64) (Rand, error) {
	u := &Uniform{}
	_ = u.init(seed, nRange, center)

	return u, nil
}
