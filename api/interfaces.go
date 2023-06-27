package api

import "time"

type SequenceGenerator interface {
	Next() float64
	NextFor(interval time.Duration) (bool, *float64, int64)
	Size() int64
	GetStartTimestamp(interval time.Duration) int64
	AdjustTime(timestamp int64)
}
