package pattern

type Generator struct {
	pattern IOPattern

	pageOffset int64
	pageSz     int64
}

func (r *Generator) Offset() int64 {
	return r.pageOffset + (r.pattern.PageNo() * r.pageSz)
}
