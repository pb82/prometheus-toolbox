package metrics

import "github.com/prometheus/client_golang/prometheus"

var Registry = prometheus.NewRegistry()

var (
	RemoteWriteRequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "prometheus_toolbox",
		Subsystem: "remote_write",
		Name:      "request_count",
		Help:      "number of received remote write requests",
	}, []string{"origin"})

	RemoteWriteRequestCompressedSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "prometheus_toolbox",
		Subsystem: "remote_write",
		Name:      "compressed_size",
		Help:      "compressed remote write request size in bytes",
	}, []string{"origin"})

	RemoteWriteRequestUncompressedSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "prometheus_toolbox",
		Subsystem: "remote_write",
		Name:      "uncompressed_size",
		Help:      "uncompressed remote write request size in bytes",
	}, []string{"origin"})

	RemoteWriteRequestTimeseriesCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "prometheus_toolbox",
		Subsystem: "remote_write",
		Name:      "timeseries_count",
		Help:      "number of time series in remote write request",
	}, []string{"origin"})

	RemoteWriteHeader = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "prometheus_toolbox",
		Subsystem: "remote_write",
		Name:      "header",
		Help:      "present http headers in remote write requests",
	}, []string{"origin", "header", "value"})
)

func init() {
	Registry.MustRegister(RemoteWriteRequestCount)
	Registry.MustRegister(RemoteWriteRequestCompressedSize)
	Registry.MustRegister(RemoteWriteRequestUncompressedSize)
	Registry.MustRegister(RemoteWriteRequestTimeseriesCount)
	Registry.MustRegister(RemoteWriteHeader)
}
