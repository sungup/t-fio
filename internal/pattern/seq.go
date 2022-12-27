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

type SeqOptions struct {
	StartFrom float64
}

func (s *SeqOptions) MakeIOPattern(nRange int64) (IOPattern, error) {
	return &Sequencer{
		until:  nRange,
		cursor: int64(float64(nRange)*s.StartFrom) % nRange,
	}, nil
}
