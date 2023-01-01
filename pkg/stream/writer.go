package stream

import (
	"context"
	"fmt"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/remotewrite"
	"github.com/pb82/prometheus-toolbox/pkg/sequence"
	"github.com/pb82/prometheus-toolbox/pkg/timeseries"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"log"
	"net/url"
	"sync"
	"time"
)

func StartStreamWriters(ctx context.Context, config *api.Config, prometheusUrl *url.URL, wg *sync.WaitGroup) error {
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
					err := remotewrite.SendWriteRequest(wr, prometheusUrl)
					if err != nil {
						log.Println(fmt.Sprintf("error sending sample: %v", err.Error()))
					}
				}
			}
		}()
	}
	return nil
}
