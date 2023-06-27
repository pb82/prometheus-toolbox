package api

import "time"

type Stream struct {
	timesAlready int64

	Initial   float64
	Increment float64
}

// NextFor not implemented for this type
func (s *Stream) NextFor(interval time.Duration) (bool, *float64, int64) {
	panic("NextFor not implemented for this type")
}

// Size not implemented for this type
func (s *Stream) Size() int64 {
	panic("Size not implemented for this type")
}

// GetStartTimestamp not implemented for this type
func (s *Stream) GetStartTimestamp(interval time.Duration) int64 {
	panic("GetStartTimestamp not implemented for this type")
}

// AdjustTime not implemented for this type
func (s *Stream) AdjustTime(timestamp int64) {
	panic("AdjustTime not implemented for this type")
}

func (s *Stream) Next() float64 {
	next := s.Initial + (float64(s.timesAlready) * s.Increment)
	s.timesAlready += 1
	return next
}
