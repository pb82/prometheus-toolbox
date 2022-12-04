package remotewrite

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/snappy"
)

func Precheck(wr *prometheus.WriteRequest, i string) error {
	now := time.Now().UnixMilli()
	interval, err := time.ParseDuration(i)
	if err != nil {
		return err
	}

	for _, ts := range wr.Timeseries {
		var lastTimestamp int64
		for _, sample := range ts.Samples {
			if sample.Timestamp > now {
				return errors.New("future sample detected")
			}

			diff := sample.Timestamp - lastTimestamp
			if diff < interval.Milliseconds() {
				return errors.New("duplicate sample detected")
			}
			lastTimestamp = sample.Timestamp
		}
	}
	return nil
}

func SendWriteRequest(wr *prometheus.WriteRequest, prometheusUrl *url.URL) error {
	data, _ := proto.Marshal(wr)
	encoded := snappy.Encode(nil, data)

	body := bytes.NewReader(encoded)
	req, err := http.NewRequest("POST", prometheusUrl.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Do(req.WithContext(context.TODO()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == 400 {
			// possibly duplicate data? we're not concerned about that
			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.New(fmt.Sprintf("unexpected remote write status code %v", resp.StatusCode))
			}
			return errors.New(fmt.Sprintf("invalid remote write request: %v", string(bytes)))
		}

		return errors.New(fmt.Sprintf("unexpected remote write status code %v", resp.StatusCode))
	}

	return nil
}
