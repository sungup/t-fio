package pattern

type Sequencer struct {
	until  int64
	cursor int64
}

func (s *Sequencer) PageNo() int64 {
	pageNo := s.cursor
	if s.cursor++; s.cursor == s.until {
		s.cursor = 0
	}

	return pageNo
}

func NewSequencer(max int64) IOPattern {

	return &Sequencer{
		until:  max,
		cursor: 0,
	}
}
