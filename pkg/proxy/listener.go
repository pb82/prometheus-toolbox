package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/pb82/prometheus-toolbox/pkg/metrics"
	"github.com/pb82/prometheus-toolbox/pkg/remotewrite"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
)

const (
	HeaderContentLength = "Content-Length"
	MaxLabelKeyLength   = 127
	MaxLabelValueLength = 255
)

type router struct {
	prometheusHandler http.Handler
}

func newRouter() router {
	return router{
		prometheusHandler: promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}),
	}
}

func handleRemoteWriteRequest(w http.ResponseWriter, r *http.Request) {
	trimToMaxSize := func(maxLength int, original string) string {
		if len(original) >= maxLength {
			return fmt.Sprintf("%v...", original[:maxLength-3])
		} else {
			return original
		}
	}

	si, _, err := remotewrite.DecodeWriteRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("received remote write request from %v: %v", r.RemoteAddr, si.String())
	metrics.RemoteWriteRequestCount.WithLabelValues(r.RemoteAddr).Inc()
	metrics.RemoteWriteRequestCompressedSize.WithLabelValues(r.RemoteAddr).Set(si.CompressedSize)
	metrics.RemoteWriteRequestUncompressedSize.WithLabelValues(r.RemoteAddr).Set(si.UncompressedSize)
	metrics.RemoteWriteRequestTimeseriesCount.WithLabelValues(r.RemoteAddr).Set(si.TimeseriesCount)

	for header, value := range r.Header {
		// the value of Content-Length is by its nature variable and would produce a large number
		// of time series if recorded
		if header == HeaderContentLength {
			continue
		}
		metrics.RemoteWriteHeader.WithLabelValues(r.RemoteAddr,
			trimToMaxSize(MaxLabelKeyLength, header),
			trimToMaxSize(MaxLabelValueLength, strings.Join(value, ","))).Set(1)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (r router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch fmt.Sprintf("%v %v", req.Method, req.URL.Path) {
	case "GET /":
		r.prometheusHandler.ServeHTTP(resp, req)
	case "POST /":
		handleRemoteWriteRequest(resp, req)
	default:
		resp.WriteHeader(http.StatusNotFound)
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
