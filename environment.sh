#!/usr/bin/env sh

PROMETHEUS_NAME="prometheus_toolbox_prometheus"
PROMETHEUS_IMAGE="docker.io/prom/prometheus:latest"
PROMETHEUS_PORT=9090

GRAFANA_NAME="prometheus_toolbox_grafana"
GRAFANA_IMAGE="docker.io/grafana/grafana-oss:latest"
GRAFANA_PORT=3000
GRAFANA_DATASOURCE=""

RUNTIME="docker"
OPTS=""


if command -v podman &> /dev/null
then
  RUNTIME=podman
  OPTS="--security-opt label=disable"
fi

if ! command -v $RUNTIME &> /dev/null
then
  echo "neither podman nor docker were found, aborting"
  exit 1
fi

echo "detected runtime is $RUNTIME, starting containers"

$RUNTIME run -d --rm --name $PROMETHEUS_NAME --network host $OPTS -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml $PROMETHEUS_IMAGE --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml &> /dev/null
PROMETHEUS_CONTAINER=`$RUNTIME ps -f name=$PROMETHEUS_NAME --format "{{.ID}}"`

$RUNTIME run -d --rm --name $GRAFANA_NAME --network host $OPTS $GRAFANA_IMAGE &> /dev/null
GRAFANA_CONTAINER=`$RUNTIME ps -f name=$GRAFANA_NAME --format "{{.ID}}"`

echo "creating prometheus datasource..."
curl --output /dev/null --silent --retry 5 --retry-connrefused --retry-delay 1 -H "Content-Type: application/json;charset=UTF-8" -XPOST http://admin:admin@localhost:$GRAFANA_PORT/api/datasources --data-binary '{"name":"prometheus","type":"prometheus","url":"http://localhost:9090","access":"proxy"}'

cleanup() {
  echo "stopping prometheus container $PROMETHEUS_CONTAINER"
  $RUNTIME stop $PROMETHEUS_CONTAINER &> /dev/null
  echo "stopping grafana container $GRAFANA_CONTAINER"
  $RUNTIME stop $GRAFANA_CONTAINER &> /dev/null
  echo "done"
  exit 0
}

trap 'cleanup' SIGINT
echo "Prometheus is running on port $PROMETHEUS_PORT"
echo "Grafana is running on port $GRAFANA_PORT"
echo "Grafana credentials are admin/admin"
echo "press ctrl-c to quit"
while true; do
  sleep 1
done