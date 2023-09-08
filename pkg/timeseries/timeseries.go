package timeseries

import (
	prometheus "buf.build/gen/go/prometheus/prometheus/protocolbuffers/go"
)

func ScanAndParseTimeSeries(source string) (*prometheus.TimeSeries, error) {
	scanner := NewScanner(source)
	scanner.Scan()

	parser := NewParser(scanner.Tokens)
	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &parser.Series, nil
}
