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

func (s *SequenceList) Next(interval time.Duration) (bool, *float64, int64) {
	if s.index >= len(s.sequences) {
		return false, nil, 0
	}
	valid, next := s.sequences[s.index].Next()
	if !valid {
		s.index += 1
		return s.Next(interval)
	}

	ts := s.startTimestamp + (s.timesAlready * interval.Milliseconds())
	s.timesAlready += 1
	return true, next, ts
}

func (s *SequenceList) AsIntArray(interval time.Duration) []int {
	var result []int
	for true {
		valid, next, _ := s.Next(interval)
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

func (s *SequenceList) GetStartTimestamp(interval time.Duration) int64 {
	return time.Now().UnixMilli() - (s.Size() * interval.Milliseconds())
}

// AdjustTime set the clock to the given timestamp
func (s *SequenceList) AdjustTime(timestamp int64) {
	s.startTimestamp = timestamp
}

// Append append a new sequence to the list
func (s *SequenceList) Append(sequence Sequence) {
	s.sequences = append(s.sequences, sequence)
}
