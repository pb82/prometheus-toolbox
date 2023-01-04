#!/usr/bin/env bash

PROMETHEUS_NAME="prometheus_toolbox_prometheus"
IMAGE="docker.io/prom/prometheus:latest"
RUNTIME="docker"
OPTS=""
PORT=9090

stty -echoctl

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

echo "detected runtime is $RUNTIME"

$RUNTIME run -d --rm --name $PROMETHEUS_NAME -p $PORT:$PORT $OPTS -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml $IMAGE --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml &> /dev/null
PROMETHEUS_CONTAINER=`podman ps -f name=prometheus_toolbox_prometheus --format "{{.ID}}"`

cleanup() {
  echo "stopping prometheus container $PROMETHEUS_CONTAINER"
  $RUNTIME stop $PROMETHEUS_CONTAINER &> /dev/null
  echo "done"
  exit 0
}

trap 'cleanup' SIGINT
echo "Prometheus is running on port $PORT, press ctrl-c to quit"
while true; do
  sleep 1
done