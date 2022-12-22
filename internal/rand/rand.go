package rand

type Randomizer interface {
	PageNo() int64
}

type Rand struct {
	rnd Randomizer

	pageSz int64
}

func (r *Rand) Offset() int64 {
	return r.rnd.PageNo() * r.pageSz
}
