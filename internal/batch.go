package internal

import "go.buf.build/protocolbuffers/go/prometheus/prometheus"

type Batch struct {
	maxSize int
	samples int
	request prometheus.WriteRequest
}

func (b *Batch) Full() bool {
	return b.samples >= b.maxSize
}

func (b *Batch) HasSamples() bool {
	return b.samples > 0
}

func (b *Batch) Reset() {
	b.samples = 0
	for _, timeseries := range b.request.Timeseries {
		timeseries.Samples = []*prometheus.Sample{}
	}
}

func (b *Batch) GetWriteRequest() *prometheus.WriteRequest {
	return &b.request
}

func (b *Batch) AddSample(timeseries *prometheus.TimeSeries, value float64, timestamp int64) bool {
	if b.Full() {
		return false
	}

	b.ensureTimeSeries(timeseries)
	timeseries.Samples = append(timeseries.Samples, &prometheus.Sample{
		Value:     value,
		Timestamp: timestamp,
	})

	b.samples += 1
	return true
}

func (b *Batch) findTimeSeries(timeseries *prometheus.TimeSeries) *prometheus.TimeSeries {
	for _, targetSeries := range b.request.Timeseries {
		if targetSeries == timeseries {
			return timeseries
		}
	}
	return nil
}

func (b *Batch) ensureTimeSeries(timeseries *prometheus.TimeSeries) {
	if b.findTimeSeries(timeseries) == nil {
		b.request.Timeseries = append(b.request.Timeseries, timeseries)
	}
}

func NewBatch(size int) *Batch {
	return &Batch{
		maxSize: size,
		samples: 0,
		request: prometheus.WriteRequest{},
	}
}
