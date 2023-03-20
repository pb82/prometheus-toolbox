#!/usr/bin/env sh

PROMETHEUS_NAME="prometheus_toolbox_prometheus"
: "${PROMETHEUS_IMAGE:="docker.io/prom/prometheus:latest"}"
PROMETHEUS_PORT=9090

GRAFANA_NAME="prometheus_toolbox_grafana"
: "${GRAFANA_IMAGE:="docker.io/grafana/grafana-oss:latest"}"
GRAFANA_PORT=3000

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

echo "detected runtime is $RUNTIME, starting containers. this can take some time..."

if [ -f ./rules.yml ]; then
  echo "prometheus rules file found, importing alerts"
  # Creata a basic prometheus config file with alerting
  cat > prometheus.yml <<- EOF
  global:
    scrape_interval: 10s
    evaluation_interval: 30s
  rule_files:
    - ./rules.yml
EOF
  # start the prometheus container with mounted rules
  echo "starting prometheus container from image $PROMETHEUS_IMAGE"
  $RUNTIME run -d --rm --name $PROMETHEUS_NAME -p 9090:9090 $OPTS -v $(pwd)/rules.yml:/etc/prometheus/rules.yml -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml $PROMETHEUS_IMAGE --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml &> /dev/null
else
  cat > prometheus.yml <<- EOF
  global:
    scrape_interval: 10s
    evaluation_interval: 30s
EOF
  # start the prometheus container
  echo "starting prometheus container from image $PROMETHEUS_IMAGE"
  $RUNTIME run -d --rm --name $PROMETHEUS_NAME -p 9090:9090 $OPTS -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml $PROMETHEUS_IMAGE --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml &> /dev/null
fi

PROMETHEUS_CONTAINER=`$RUNTIME ps -f name=$PROMETHEUS_NAME --format "{{.ID}}"`

# start the grafana container
echo "starting grafana container from image $GRAFANA_IMAGE"
$RUNTIME run -d --rm --name $GRAFANA_NAME -p 3000:3000 $OPTS $GRAFANA_IMAGE &> /dev/null
GRAFANA_CONTAINER=`$RUNTIME ps -f name=$GRAFANA_NAME --format "{{.ID}}"`

# connect grafana and prometheus
echo "creating prometheus datasource"
curl --output /dev/null --retry 5 --retry-all-errors --retry-delay 5 --silent -H "Content-Type: application/json;charset=UTF-8" -XPOST http://admin:admin@localhost:3000/api/datasources --data-binary '{"name":"prometheus","type":"prometheus","url":"http://172.17.0.1:9090","access":"proxy"}'

cleanup() {
  echo "stopping prometheus container $PROMETHEUS_CONTAINER"
  $RUNTIME stop $PROMETHEUS_CONTAINER &> /dev/null
  echo "stopping grafana container $GRAFANA_CONTAINER"
  $RUNTIME stop $GRAFANA_CONTAINER &> /dev/null
  rm -f prometheus.yml
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
