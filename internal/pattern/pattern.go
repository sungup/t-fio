package pattern

type IOPattern interface {
	PageNo() int64
}

type Pattern struct {
	rnd IOPattern

	pageOffset int64
	pageSz     int64
}

func (r *Pattern) Offset() int64 {
	return r.pageOffset + (r.rnd.PageNo() * r.pageSz)
}
