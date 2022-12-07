package stream

import (
	"fmt"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"log"
	"net/url"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/remotewrite"
	"github.com/pb82/prometheus-toolbox/pkg/sequence"
	"github.com/pb82/prometheus-toolbox/pkg/timeseries"
	"sync"
	"time"
)

func StartStreamWriters(config *api.Config, prometheusUrl *url.URL, wg *sync.WaitGroup, stop <-chan bool) (int, error) {
	count := 0
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return count, err
	}

	for _, ts := range config.Series {
		if ts.Stream == "" {
			continue
		}

		series, err := timeseries.ScanAndParseTimeSeries(ts.Series)
		if err != nil {
			return count, err
		}

		stream, err := sequence.ScanAndParseStream(ts.Stream)
		if err != nil {
			return count, err
		}

		wg.Add(1)
		count += 1
		go func() {
			for {
				select {
				case _ = <-stop:
					wg.Done()
					return
				case <-time.After(interval):
					nextValue := stream.Next()
					timeseries := prometheus.TimeSeries{}
					timeseries.Labels = series.Labels
					timeseries.Samples = append(timeseries.Samples, &prometheus.Sample{
						Value:     nextValue,
						Timestamp: time.Now().UnixMilli(),
					})
					wr := &prometheus.WriteRequest{}
					wr.Timeseries = append(wr.Timeseries, &timeseries)
					log.Println(fmt.Sprintf("sending sample for timeseries %v: %v", ts.Series, nextValue))
					err := remotewrite.SendWriteRequest(wr, prometheusUrl)
					if err != nil {
						log.Println(fmt.Sprintf("error sending sample: %v", err.Error()))
					}
				}
			}
		}()
	}
	return count, nil
}
