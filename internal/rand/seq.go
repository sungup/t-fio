package rand

type Sequencer struct {
	base   int64
	until  int64
	cursor int64
}

func (s *Sequencer) PageNo() int64 {
	pageNo := s.cursor
	if s.cursor++; s.cursor == s.until {
		s.cursor = s.base
	}

	return pageNo
}

func NewSequencer(pageSz, baseOffset, ioRangeSz int64) Randomizer {
	base := baseOffset / pageSz

	return &Sequencer{
		base:   base,
		until:  base + ioRangeSz/pageSz,
		cursor: base,
	}
}
