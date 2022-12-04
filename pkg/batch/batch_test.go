package batch

import (
	"promtoolbox/api"
	sequence2 "promtoolbox/pkg/sequence"
	"promtoolbox/pkg/timeseries"
	"testing"
)

func TestBatch(t *testing.T) {
	type testcase struct {
		inputConfig          api.Config
		wantNumberOfRequests int
		wantNumberOfSamples  map[int]int
		batchSize            int
	}

	testcases := []testcase{
		{
			inputConfig: api.Config{
				Series: []api.TimeseriesConfig{
					{
						Series: "up",
						Values: "1+0x9",
					},
					{
						Series: "down",
						Values: "1+0x9",
					},
				},
			},
			wantNumberOfRequests: 4,
			wantNumberOfSamples: map[int]int{
				0: 5,
				1: 4,
				2: 5,
				3: 4,
			},
			batchSize: 5,
		},

		{
			inputConfig: api.Config{
				Series: []api.TimeseriesConfig{
					{
						Series: "up",
						Values: "1+0x9",
					},
				},
			},
			wantNumberOfRequests: 2,
			wantNumberOfSamples: map[int]int{
				0: 5,
				1: 4,
			},
			batchSize: 5,
		},
		{
			inputConfig: api.Config{
				Series: []api.TimeseriesConfig{
					{
						Series: "up",
						Values: "1+0x10",
					},
				},
			},
			wantNumberOfRequests: 1,
			wantNumberOfSamples: map[int]int{
				0: 10,
			},
			batchSize: 10,
		},
		{
			inputConfig: api.Config{
				Series: []api.TimeseriesConfig{
					{
						Series: "up",
						Values: "1+0x100",
					},
				},
			},
			wantNumberOfRequests: 10,
			wantNumberOfSamples: map[int]int{
				0: 10,
				1: 10,
				2: 10,
				3: 10,
				4: 10,
				5: 10,
				6: 10,
				7: 10,
				8: 10,
				9: 10,
			},
			batchSize: 10,
		},
		{
			inputConfig: api.Config{
				Series: []api.TimeseriesConfig{
					{
						Series: "up",
						Values: "1+0x99",
					},
				},
			},
			wantNumberOfRequests: 10,
			wantNumberOfSamples: map[int]int{
				0: 10,
				1: 10,
				2: 10,
				3: 10,
				4: 10,
				5: 10,
				6: 10,
				7: 10,
				8: 10,
				9: 9,
			},
			batchSize: 10,
		},
	}

	for _, tc := range testcases {
		batch := NewBatch(tc.batchSize)

		for _, ts := range tc.inputConfig.Series {
			series, err := timeseries.ScanAndParseTimeSeries(ts.Series)
			if err != nil {
				t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
			}

			sequence, err := sequence2.ScanAndParseSequence(ts.Values)
			if err != nil {
				t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
			}

			for true {
				valid, value, timestamp := sequence.Next()
				if !valid {
					break
				}

				if value == nil {
					continue
				}

				batch.AddSample(series, timestamp, *value)
			}
		}

		writeRequests := batch.GetWriteRequests()

		if len(writeRequests) != tc.wantNumberOfRequests {
			t.Fatalf("\nwant requests: %v \n got requests: %v", tc.wantNumberOfRequests, len(writeRequests))
		}

		for k, v := range tc.wantNumberOfSamples {
			if len(writeRequests[k].Timeseries[0].Samples) != v {
				t.Fatalf("\nwant samples: %v \n got samples: %v", v, len(writeRequests[k].Timeseries[0].Samples))
			}
		}
	}
}
