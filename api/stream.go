package api

type Stream struct {
	timesAlready int64

	Initial   float64
	Increment float64
}

func (s *Stream) Next() float64 {
	next := s.Initial + (float64(s.timesAlready) * s.Increment)
	s.timesAlready += 1
	return next
}
