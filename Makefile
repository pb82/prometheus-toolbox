build:
	go build -o promtoolbox cmd/main.go

prom:
	docker run -p 9090:9090 --network=host -v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml