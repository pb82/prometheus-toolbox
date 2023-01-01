package internal

import "fmt"

// SizeInfo contains information about the size of a remote write request
type SizeInfo struct {
	CompressedSize   float64
	UncompressedSize float64
	TimeseriesCount  float64
}

func (s *SizeInfo) String() string {
	return fmt.Sprintf("compressed: %v bytes, uncompressed: %v bytes, times series count: %v",
		s.CompressedSize, s.UncompressedSize, s.TimeseriesCount)
}
