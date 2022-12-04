package precalculated

import (
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"promtoolbox/api"
	batch2 "promtoolbox/pkg/batch"
	sequence2 "promtoolbox/pkg/sequence"
	"promtoolbox/pkg/timeseries"
	"time"
)

func GetPrecalculatedRemoteWriteRequests(config *api.Config, batchSize int) ([]prometheus.WriteRequest, error) {
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return nil, err
	}

	batch := batch2.NewBatch(batchSize)

	for _, ts := range config.Series {
		timeSeries, err := timeseries.ScanAndParseTimeSeries(ts.Series)
		if err != nil {
			return nil, err
		}

		sequence, err := sequence2.ScanAndParseSequence(ts.Values)
		if err != nil {
			return nil, err
		}
		sequence.AdjustTime(interval)

		samples := 0

		for true {
			if samples >= batchSize {

			}

			valid, value, timestamp := sequence.Next()
			if !valid {
				break
			}

			batch.AddSample(timeSeries, *value, timestamp)
		}
	}

	return batch.GetWriteRequests(), nil
}
