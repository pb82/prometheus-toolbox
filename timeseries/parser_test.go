package timeseries

import (
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	type testcase struct {
		input          string
		wantErr        bool
		wantTimeSeries prometheus.TimeSeries
	}

	testcases := []testcase{
		{
			input:   "metric{label=\"value\"}",
			wantErr: false,
			wantTimeSeries: prometheus.TimeSeries{
				Labels: []*prometheus.Label{
					{
						Name:  "__name__",
						Value: "metric",
					},
					{
						Name:  "label",
						Value: "value",
					},
				},
			},
		},
		{
			input:   "metric{label1=\"value1\", label2=\"value2\"}",
			wantErr: false,
			wantTimeSeries: prometheus.TimeSeries{
				Labels: []*prometheus.Label{
					{
						Name:  "__name__",
						Value: "metric",
					},
					{
						Name:  "label1",
						Value: "value1",
					},
					{
						Name:  "label2",
						Value: "value2",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		scanner := NewScanner(tc.input)
		scanner.Scan()

		parser := NewParser(scanner.Tokens)
		err := parser.Parse()
		if err != nil && !tc.wantErr {
			t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
		}
		if !reflect.DeepEqual(parser.Series, tc.wantTimeSeries) {
			t.Fatalf("\nwant: %v \n got: %v", tc.wantTimeSeries, parser.Series)
		}
	}

}
