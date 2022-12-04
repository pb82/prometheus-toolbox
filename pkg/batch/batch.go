package batch

import (
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
)

type batchSample struct {
	series    *prometheus.TimeSeries
	timestamp int64
	value     float64
}

type Batch struct {
	samplesPerBatch int
	batches         map[*prometheus.TimeSeries][]*prometheus.Sample
	samples         []batchSample
	series          map[*prometheus.TimeSeries]bool
}

// AddSample add a sample for the given time series at the given timestamp to the batch
func (b *Batch) AddSample(ts *prometheus.TimeSeries, timestamp int64, value float64) {
	b.samples = append(b.samples, batchSample{
		series:    ts,
		timestamp: timestamp,
		value:     value,
	})
	b.series[ts] = true
	/*
		b.batches[ts] = append(b.batches[ts], &prometheus.Sample{
			Value:     value,
			Timestamp: timestamp,
		})
	*/
}

// GetWriteRequests returns the list of remote write requests required to fit all samples given a batch size
func (b *Batch) GetWriteRequests() []*prometheus.WriteRequest {
	numWriteRequests := len(b.samples) / b.samplesPerBatch
	if numWriteRequests <= 0 {
		numWriteRequests = 1
	}

	writeRequests := make([]*prometheus.WriteRequest, numWriteRequests)
	i := 0
	for i < numWriteRequests {
		writeRequests[i] = new(prometheus.WriteRequest)
		i += 1
	}

	requestMapping := map[*prometheus.TimeSeries][]*prometheus.TimeSeries{}

	for _, writeRequest := range writeRequests {
		for series, _ := range b.series {
			ts := &prometheus.TimeSeries{
				Labels:  series.Labels,
				Samples: []*prometheus.Sample{},
			}
			writeRequest.Timeseries = append(writeRequest.Timeseries, ts)
			requestMapping[series] = append(requestMapping[series], ts)
		}
	}

	batchIndex := -1
	for index, sample := range b.samples {
		if index%b.samplesPerBatch == 0 {
			batchIndex += 1
		}
		ts := requestMapping[sample.series][batchIndex]
		ts.Samples = append(ts.Samples, &prometheus.Sample{
			Value:     sample.value,
			Timestamp: sample.timestamp,
		})
	}

	/*
		for timeseries, samples := range b.batches {
			iterations := 0
			for i := 0; i < len(samples); i += b.samplesPerBatch {
				limit := b.samplesPerBatch
				samplesLeft := len(samples) - (iterations * b.samplesPerBatch)
				if limit > samplesLeft {
					limit = samplesLeft
				}
				writeRequests = append(writeRequests, prometheus.WriteRequest{
					Timeseries: []*prometheus.TimeSeries{
						{
							Labels:  timeseries.Labels,
							Samples: samples[i : i+limit],
						},
					},
					Metadata: nil,
				})
				iterations += 1
			}
		}
	*/
	return writeRequests
}

func NewBatch(size int) *Batch {
	return &Batch{
		samplesPerBatch: size,
		batches:         map[*prometheus.TimeSeries][]*prometheus.Sample{},
		series:          map[*prometheus.TimeSeries]bool{},
		samples:         []batchSample{},
	}
}
