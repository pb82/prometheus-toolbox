package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/snappy"
	"github.com/pb82/prometheus-toolbox/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
)

type sizeInfo struct {
	compressedSize   float64
	uncompressedSize float64
}

type router struct {
	prometheusHandler http.Handler
}

func newRouter() router {
	return router{
		prometheusHandler: promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}),
	}
}

func decodeWriteRequest(r io.Reader) (*sizeInfo, *prometheus.WriteRequest, error) {
	var si sizeInfo
	compressed, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, nil, err
	}

	si.compressedSize = float64(len(compressed))
	si.uncompressedSize = float64(len(reqBuf))

	var req prometheus.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, nil, err
	}

	return &si, &req, nil
}

func handleRemoteWriteRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("remote write request received")

	si, remoteWriteRequest, err := decodeWriteRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.RemoteWriteRequestCount.WithLabelValues(r.RemoteAddr).Inc()
	metrics.RemoteWriteRequestCompressedSize.WithLabelValues(r.RemoteAddr).Set(si.compressedSize)
	metrics.RemoteWriteRequestUncompressedSize.WithLabelValues(r.RemoteAddr).Set(si.uncompressedSize)
	metrics.RemoteWriteRequestTimeseriesCount.WithLabelValues(r.RemoteAddr).Set(float64(len(remoteWriteRequest.Timeseries)))

	w.WriteHeader(http.StatusNoContent)
}

func (r router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == "/" {
		handleRemoteWriteRequest(resp, req)
	} else if req.Method == http.MethodGet && req.URL.Path == "/" {
		r.prometheusHandler.ServeHTTP(resp, req)
	}
}

func StartListener(ctx context.Context, port *int) error {
	if port == nil {
		return errors.New("port number required")
	}

	if *port <= 0 || *port > 65536 {
		return errors.New(fmt.Sprintf("%v is not a valid port number", *port))
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", *port),
		Handler: newRouter(),
	}

	go func() {
		log.Printf("listening for remote write POST requests on port %v", *port)
		log.Println("use GET on the same port to view remote write metrics")

		server.ListenAndServe()
	}()

	<-ctx.Done()
	server.Shutdown(ctx)
	return nil
}
