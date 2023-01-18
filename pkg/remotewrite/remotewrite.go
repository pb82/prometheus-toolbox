package remotewrite

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/pb82/prometheus-toolbox/internal"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	clientTimeout = 30 * time.Second
)

var (
	httpTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient = http.Client{
		Timeout:   clientTimeout,
		Transport: httpTransport,
	}
)

type RemoteWriter struct {
	client        http.Client
	prometheusUrl *url.URL
}

func NewRemoteWriterWithOIDCTransport(ctx context.Context, prometheusUrl *url.URL, oidcConfig *internal.OIDCConfig) (*RemoteWriter, error) {
	provider, err := oidc.NewProvider(ctx, oidcConfig.IssuerUrl)
	if err != nil {
		return nil, err
	}

	var cfg = clientcredentials.Config{
		ClientID:     oidcConfig.ClientId,
		ClientSecret: oidcConfig.ClientSecret,
		TokenURL:     provider.Endpoint().TokenURL,
	}

	if oidcConfig.Audience != "" {
		cfg.EndpointParams = map[string][]string{
			"audience": {oidcConfig.Audience},
		}
	}

	return &RemoteWriter{
		client: http.Client{
			Transport: &oauth2.Transport{
				Source: cfg.TokenSource(ctx),
				Base:   httpTransport,
			},
			Timeout: clientTimeout,
		},
		prometheusUrl: prometheusUrl,
	}, nil
}

func NewRemoteWriter(prometheusUrl *url.URL) *RemoteWriter {
	return &RemoteWriter{
		client:        httpClient,
		prometheusUrl: prometheusUrl,
	}
}

// SendWriteRequest encodes and sends a write request to the given url
func (s *RemoteWriter) SendWriteRequest(wr *prometheus.WriteRequest) error {
	data, _ := proto.Marshal(wr)
	encoded := snappy.Encode(nil, data)

	body := bytes.NewReader(encoded)
	req, err := http.NewRequest("POST", s.prometheusUrl.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	resp, err := s.client.Do(req)
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

// DecodeWriteRequest deserialize a compressed remote write request
func DecodeWriteRequest(r io.Reader) (*internal.SizeInfo, *prometheus.WriteRequest, error) {
	var si internal.SizeInfo

	compressed, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, nil, err
	}

	si.CompressedSize = float64(len(compressed))
	si.UncompressedSize = float64(len(reqBuf))

	var req prometheus.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, nil, err
	}

	si.TimeseriesCount = float64(len(req.Timeseries))
	return &si, &req, nil
}
