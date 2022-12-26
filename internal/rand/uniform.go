package rand

type Uniform struct {
	core
}

// Uint64 is a random generating function but, there is no need hashing because all generated
// random values has same probability.
func (u *Uniform) Uint64() uint64 {
	return u.core.rand.Uint64() % u.nRange
}

func NewUniform(seed int64, nRange uint64) (Rand, error) {
	u := &Uniform{}
	_ = u.init(seed, nRange, 0)

	return u, nil
}
