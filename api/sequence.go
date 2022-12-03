package api

import "time"

type Sequence struct {
	Initial      float64
	Increment    float64
	Times        int64
	TimesAlready int64
}

type SequenceList struct {
	index          int
	sequences      []Sequence
	startTimestamp int64
	timesAlready   int64
	interval       time.Duration
}

func (s *Sequence) Next() (bool, *float64) {
	if s.TimesAlready >= s.Times {
		return false, nil
	}

	nextValue := s.Initial + (float64(s.TimesAlready) * s.Increment)
	s.TimesAlready += 1
	return true, &nextValue
}

func (s *Sequence) Size() int64 {
	return s.Times
}

func (s *SequenceList) Next() (bool, *float64, int64) {
	if s.index >= len(s.sequences) {
		return false, nil, 0
	}
	valid, next := s.sequences[s.index].Next()
	if !valid {
		s.index += 1
		return s.Next()
	}

	ts := s.startTimestamp + (s.timesAlready * s.interval.Milliseconds())
	s.timesAlready += 1
	return true, next, ts
}

func (s *SequenceList) AsIntArray() []int {
	var result []int
	for true {
		valid, next, _ := s.Next()
		if !valid {
			break
		}
		result = append(result, int(*next))
	}
	return result
}

// Size returns the number of iterations over all sequences in the list
func (s *SequenceList) Size() int64 {
	var size int64 = 0
	for _, seq := range s.sequences {
		size += seq.Size()
	}
	return size
}

// AdjustTime rewind the clock so that the series goes back far enough to fit all samples
func (s *SequenceList) AdjustTime(interval time.Duration) {
	s.interval = interval
	s.startTimestamp = time.Now().UnixMilli() - (s.Size() * interval.Milliseconds())
}

// Append append a new sequence to the list
func (s *SequenceList) Append(sequence Sequence) {
	s.sequences = append(s.sequences, sequence)
}
