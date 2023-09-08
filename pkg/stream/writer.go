package stream

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	prometheus "buf.build/gen/go/prometheus/prometheus/protocolbuffers/go"

	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/remotewrite"
	"github.com/pb82/prometheus-toolbox/pkg/sequence"
	"github.com/pb82/prometheus-toolbox/pkg/timeseries"
)

func StartStreamWriters(ctx context.Context, config *api.Config, rw *remotewrite.RemoteWriter, wg *sync.WaitGroup) error {
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return err
	}

	for _, ts := range config.Series {
		ts := ts
		if ts.Stream == "" || ts.Series == "" {
			continue
		}

		series, err := timeseries.ScanAndParseTimeSeries(ts.Series)
		if err != nil {
			return err
		}

		stream, err := sequence.ScanAndParseStream(ts.Stream)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case <-time.After(interval):
					nextValue := stream.Next()
					sendSeries := prometheus.TimeSeries{}
					sendSeries.Labels = series.Labels
					sendSeries.Samples = append(sendSeries.Samples, &prometheus.Sample{
						Value:     nextValue,
						Timestamp: time.Now().UnixMilli(),
					})
					wr := &prometheus.WriteRequest{}
					wr.Timeseries = append(wr.Timeseries, &sendSeries)
					log.Println(fmt.Sprintf("sending sample for timeseries %v: %v", ts.Series, nextValue))
					err := rw.SendWriteRequest(wr)
					if err != nil {
						log.Println(fmt.Sprintf("error sending sample: %v", err.Error()))
					}
				}
			}
		}()
	}
	return nil
}
