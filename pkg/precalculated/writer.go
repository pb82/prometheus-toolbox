package precalculated

import (
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"promtoolbox/api"
	"time"
)

func GetPrecalculatedRemoteWriteRequests(config *api.Config) ([]*prometheus.WriteRequest, error) {
	var writeRequests []*prometheus.WriteRequest
	_, err := time.ParseDuration(config.Interval)
	if err != nil {
		return nil, err
	}

	return writeRequests, nil
}