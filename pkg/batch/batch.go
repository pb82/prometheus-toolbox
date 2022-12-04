package batch

import "go.buf.build/protocolbuffers/go/prometheus/prometheus"

type Batch struct {
	maxSize  int
	samples  int
	request  prometheus.WriteRequest
	requests []prometheus.WriteRequest
}

func (b *Batch) full() bool {
	return b.samples >= b.maxSize
}

func (b *Batch) reset() {
	b.samples = 0
	b.request = prometheus.WriteRequest{}
}

func (b *Batch) GetWriteRequests() []prometheus.WriteRequest {
	if b.samples > 0 {
		b.requests = append(b.requests, b.request)
	}
	return b.requests
}

func (b *Batch) AddSample(timeseries *prometheus.TimeSeries, value float64, timestamp int64) {
	if b.full() {
		b.requests = append(b.requests, b.request)
		b.reset()
	}

	b.ensureTimeSeries(timeseries)
	timeseries.Samples = append(timeseries.Samples, &prometheus.Sample{
		Value:     value,
		Timestamp: timestamp,
	})

	b.samples += 1
}

func (b *Batch) AddTimeSeries(timeseries *prometheus.TimeSeries) {
	b.ensureTimeSeries(timeseries)
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
