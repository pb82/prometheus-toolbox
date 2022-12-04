package precalculated

import (
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"promtoolbox/api"
	sequence2 "promtoolbox/pkg/sequence"
	"promtoolbox/pkg/timeseries"
	"time"
)

// SchedulePrecalculatedRemoteWriteRequests distribute samples to remote write requests so that no single request exceeds the batch size
// and samples are intermingled to avoid out-of-order or out-of-bounds errors
func SchedulePrecalculatedRemoteWriteRequests(config *api.Config, batchSize int) ([]*prometheus.WriteRequest, int64, error) {
	type generator struct {
		ts  *prometheus.TimeSeries
		seq *api.SequenceList
	}

	interval, err := time.ParseDuration(config.Interval)
	var totalSamples int64

	if err != nil {
		return nil, totalSamples, err
	}

	var generators []generator

	for _, ts := range config.Series {
		series, err := timeseries.ScanAndParseTimeSeries(ts.Series)
		if err != nil {
			return nil, totalSamples, err
		}

		sequence, err := sequence2.ScanAndParseSequence(ts.Values)
		if err != nil {
			return nil, totalSamples, err
		}
		sequence.AdjustTime(interval)
		totalSamples += sequence.Size()
		generators = append(generators, generator{
			ts:  series,
			seq: sequence,
		})
	}

	var writeRequests []*prometheus.WriteRequest
	var currentRequest *prometheus.WriteRequest
	var scheduledSamples int64 = 0
	tsMapping := map[*prometheus.TimeSeries]*prometheus.TimeSeries{}

	for scheduledSamples < totalSamples {
		if scheduledSamples%int64(batchSize) == 0 {
			currentRequest = new(prometheus.WriteRequest)
			for _, g := range generators {
				ts := &prometheus.TimeSeries{
					Labels:  g.ts.Labels,
					Samples: []*prometheus.Sample{},
				}
				currentRequest.Timeseries = append(currentRequest.Timeseries, ts)
				tsMapping[g.ts] = ts
			}
			writeRequests = append(writeRequests, currentRequest)
		}
		scheduledSamples += 1
		g := generators[scheduledSamples%int64(len(generators))]
		valid, value, timestamp := g.seq.Next()
		if !valid {
			continue
		}
		if value == nil {
			continue
		}

		ts := tsMapping[g.ts]
		ts.Samples = append(ts.Samples, &prometheus.Sample{
			Value:     *value,
			Timestamp: timestamp,
		})
	}

	return writeRequests, totalSamples, nil
}
