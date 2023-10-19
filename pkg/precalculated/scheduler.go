package precalculated

import (
	"math"
	"time"

	prometheus "buf.build/gen/go/prometheus/prometheus/protocolbuffers/go"

	"github.com/pb82/prometheus-toolbox/api"
	sequence2 "github.com/pb82/prometheus-toolbox/pkg/sequence"
	"github.com/pb82/prometheus-toolbox/pkg/timeseries"
)

// SchedulePrecalculatedRemoteWriteRequests distribute samples to remote write requests so that no single request exceeds the batch size
// and samples are intermingled to avoid out-of-order or out-of-bounds errors
func SchedulePrecalculatedRemoteWriteRequests(config *api.Config, batchSize int) ([]*prometheus.WriteRequest, int64, error) {
	type generator struct {
		ts  *prometheus.TimeSeries
		seq api.SequenceGenerator
	}
	var totalSamples int64
	var generators []generator
	var writeRequests []*prometheus.WriteRequest
	var currentRequest *prometheus.WriteRequest
	var scheduledSamples int64 = 0

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return nil, totalSamples, err
	}

	// earliest sample timestamp over all series
	var startTimestamp int64 = math.MaxInt64

	// collect all timeseries from the config along with their sequences
	for _, ts := range config.Series {
		if ts.Series == "" || ts.Values == "" {
			continue
		}

		series, err := timeseries.ScanAndParseTimeSeries(ts.Series)
		if err != nil {
			return nil, totalSamples, err
		}

		sequence, err := sequence2.ScanAndParseSequence(ts.Values)
		if err != nil {
			return nil, totalSamples, err
		}

		start := sequence.GetStartTimestamp(interval)
		if start < startTimestamp {
			startTimestamp = start
		}

		totalSamples += sequence.Size()
		generators = append(generators, generator{
			ts:  series,
			seq: sequence,
		})
	}

	for _, g := range generators {
		g.seq.AdjustTime(startTimestamp)
	}

	// we need to split up timeseries into batches, tsMapping maps the original timeseries
	// to the one with the same labels in the current write request
	tsMapping := map[*prometheus.TimeSeries]*prometheus.TimeSeries{}

	// collect all samples
	iterations := 0
	for scheduledSamples < totalSamples {
		if iterations%batchSize == 0 {
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
		g := generators[iterations%len(generators)]
		iterations += 1

		valid, value, timestamp := g.seq.NextFor(interval)
		if !valid {
			continue
		}

		scheduledSamples += 1
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
