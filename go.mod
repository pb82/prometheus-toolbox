module prometheus-toolbox

require github.com/ghodss/yaml v1.0.0

require go.buf.build/protocolbuffers/go/prometheus/prometheus v1.3.3

require (
	github.com/golang/snappy v0.0.4
	google.golang.org/protobuf v1.28.1
)

require (
	go.buf.build/protocolbuffers/go/gogo/protobuf v1.3.9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

go 1.19
