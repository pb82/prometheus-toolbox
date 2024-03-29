# Prometheus Toolbox

A CLI tool to generate metrics data for the [Prometheus](https://prometheus.io/) monitorig system.
It allows you to define time series and values and send them directly to Prometheus using the remote write protocol.

## Why?

I found this very useful for testing PromQL queries and alerts.
Sometimes you don't have the right data in Prometheus or simply no access. 
This tool lets you simulate the data you need and test your queries against it.

## Installation

The latest version:

`GOPROXY=direct go install github.com/pb82/prometheus-toolbox@latest`

A specific version:

`GOPROXY=direct go install github.com/pb82/prometheus-toolbox@<version>`

*NOTE*: [Go](https://go.dev/) must be installed

### Building locally

Build directly for your platform with go:

```shell
$ go build
```

## Starting a local development environment

To start a local development environment including Prometheus and Grafana use the built-in `--environment` command.
This prints out a shell script to spin up containers for both applications.
It also sets up a data source in Grafana to connect it to Prometheus.
Either `podman` or `docker`, as well as `curl` are required.

Example:

```sh
$ prometheus-toolbox --environment | bash
detected runtime is podman, starting containers. this can take some time...
starting prometheus container from image docker.io/prom/prometheus:latest
starting grafana container from image docker.io/grafana/grafana-oss:latest
creating prometheus datasource
Prometheus is running on port 9090
Grafana is running on port 3000
Grafana credentials are admin/admin
press ctrl-c to quit
```

The Prometheus and Grafana images can be overridden by exporting the `PROMETHEUS_IMAGE` or `GRAFANA_IMAGE` variables:

```shell
$ export GRAFANA_IMAGE=docker.io/grafana/grafana-oss:8.5.15 && ./prometheus-toolbox --environment | bash
```

### Importing Prometheus rules into the local development environment

The environment script checks for the presence of a file with the name `rules.yml` in the directory it's running.
If present, Prometheus is configured to import [alerting](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) and [recording](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/) rules from it.

A sample rules file can be generated with the following command:

```shell
$ prometheus-toolbox --rules > rules.yml
```

# Flags

The following flags are accepted:

* `--help` Print help and exit
* `--version` Print version and exit
* `--prometheus.url` Prometheus base url
* `--config.file` Config file location, defaults to `./config.yml`
* `--batch.size` Samples per remote write request, defaults to 500
* `--proxy.listen` Turns on the remote write listener, defaults to `false`
* `--proxy.listen.port` Port to receive remote write requests, defaults to 3241
* `--environment` Print environment setup script and exit
* `--init` Print sample config file and exit
* `--rules` Print sample alerting rules file and exit
* `--oidc.enabled` Enable authenticated requests
* `--oidc.issuer` OIDC auth token issuer URL
* `--oidc.clientId` OIDC client id
* `--oidc.clientSecret` OIDC client secret
* `--oidc.audience` OIDC audience
* `--prometheus.url.suffix` Allow alternate remote write endpoints, in case your receiver uses a different endpoint than the default `/api/v1/write`

# Config file format

This tool reads from a config file where the simulated time series and values are defined.
The format is:

```yaml
interval: "10s"                       # Interval between samples, in this case 10 seconds
time_series:                          # List of time series to simulate
  - series: metric_a{label="a"}       # Time series (metric name and label list)
    values: 1+1x100                   # Precalculated samples
  - series: metric_a{label="a"}       # Another time series
    stream: 1+0                       # Realtime samples
  - series: metric_b{label="a"}       # 
    values: 1+1x50 50+0x50            # Multiple value sequences are possible
  - series: metric_c{l1="a",l2="b"}   # Multiple labels are possible
    values: _x50 1+0x50               # The underscore represents an empty value (no data received). Time still advances.    
```

The format for precalculated samples is:

```
values: 1+1x100 # <initial>(plus or minus)<increment>x<how many samples>
```

The format for realtime samples is:

```
stream: 1+1 # <initial>(plus or minus)<increment>
```

Streaming realtime samples only ends when the user interrupts the program, thus no number of samples is needed.

A basic config file can be created using the `--init` command:

```shell
$ prometheus-toolbox --init > config.yml
```

# Features

[Data generator](#data-generator): write precalculated samples.

[Streaming samples](#streaming-samples): write realtime data.

[Remote write inspector](#remote-write-inspector): inspect remote write requests.

## Data generator

Generates data and sends it to a Prometheus instance using the [remote write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) protocol.

**NOTE:** remote write needs to be specifically enabled in Prometheus with the flag `--web.enable-remote-write-receiver`

Provide the base URL of the Prometheus instance and a config file, for example:

```shell
$ cat config.yaml
---
interval: 10s
time_series:
  - series: metric_a{label="a"}
    values: 1+0x100 _x50
  - series: metric_b{label="b"}
    values: 1+1x50 50+0x50 50-1x50
```

**NOTE:** if you're not familiar with the `series` and `values` notation, it is similar to how unit tests with `promtool` are written. For more information, have a look [here](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/#series).

Run the CLI with the following options:

```shell
$ ./prometheus-toolbox --prometheus.url=http://localhost:9090 --config.file=./config.yml
```

Open the Prometheus UI and run a query like `{__name__~="metric_a|metric_b"}`, or go to the graph view.
Samples for both time series should be visible.

### How the samples are generated

Samples are always generated for the past. The time series that produces the most samples (in the above example it is `metrics_a`) determines how far back in time the samples go.
In the above case 100 samples are generated with an interval of ten seconds between them.
So both series will start at a point in time 1000 seconds in the past.

Samples are batched to avoid sending one huge remote write request.
The default batch size is 500.

## Streaming samples

You can stream continuous samples for series to a Prometheus instances:

```shell
$ cat example/config.yaml
---
interval: 10s
time_series:
  - series: up{label="a"}
    stream: 1+0
```

This will send a sample generated from the field `stream` every ten seconds.

**NOTE:** in contrast to `values`, when providing `stream` you can't specify the number of samples. Data is generated continuously.

**NOTE:** you can have both, `values` and `stream` for the same series.

## Remote write inspector

This tool can receive remote write requests and expose insights as metrics.
To enable the remote write listener, pass the `--proxy.listen` flag.
By default, remote write requests are received on port 3241, but this can be overridden using the `--proxy.listen.port` flag.

When remote write requests are received, `prometheus-toolbox` will collect data and expose it as metrics.
The metrics can be viewed on the same port as the listener with a GET request.

The following data is collected:

* Request count - how many write requests have been received
* Compressed size - snappy compressed size of the write request in bytes
* Uncompressed size - uncompressed size of the write request in bytes
* Timeseries count - number of time series in write request
* Http headers - present http headers

For example:

Configure prometheus to send write requests to `prometheus-toolbox`, by adding a `remote_write` section to `prometheus.yaml`: 

```yaml
remote_write:
  - name: "self"
    url: "http://localhost:3241"
```

**NOTE:** This is assuming that both Prometheus and prometheus-toolbox are running locally.

Then send a GET request to the `--proxy.listen` endpoint: 

```shell
$ curl -XGET http://localhost:3241

# HELP prometheus_toolbox_remote_write_compressed_size compressed remote write request size in bytes
# TYPE prometheus_toolbox_remote_write_compressed_size gauge
prometheus_toolbox_remote_write_compressed_size{origin="[::1]:58060"} 112
# HELP prometheus_toolbox_remote_write_header present http headers in remote write requests
# TYPE prometheus_toolbox_remote_write_header gauge
prometheus_toolbox_remote_write_header{header="Authorization",origin="[::1]:58060",value="Bearer 92MO6wrerqwoLWlJZfhunRIRGwl2QgsGl9rkJRI5Ts3mr1oYLBsjw92dREb2QRpfrGn5xvZFC4pwKdSupMTcBHZVRien6oXg8VSFCuAaquF9Wcagem6aHVyewDejqXFfg0m0pgTAf6hqbfNFFe5ysnygnaZfaGNHpKR5SC2VpscTKeHElEvhJUqSHxc6TO9RMzJBXjL2XPi0GMAshLMrP23clY9y9UtByF2cbobH3eYqrbRUocXog..."} 1
prometheus_toolbox_remote_write_header{header="Content-Encoding",origin="[::1]:58060",value="snappy"} 1
prometheus_toolbox_remote_write_header{header="Content-Type",origin="[::1]:58060",value="application/x-protobuf"} 1
prometheus_toolbox_remote_write_header{header="User-Agent",origin="[::1]:58060",value="Prometheus/2.39.1"} 1
prometheus_toolbox_remote_write_header{header="X-Prometheus-Remote-Write-Version",origin="[::1]:58060",value="0.1.0"} 1
# HELP prometheus_toolbox_remote_write_request_count number of received remote write requests
# TYPE prometheus_toolbox_remote_write_request_count counter
prometheus_toolbox_remote_write_request_count{origin="[::1]:58060"} 3
# HELP prometheus_toolbox_remote_write_timeseries_count number of time series in remote write request
# TYPE prometheus_toolbox_remote_write_timeseries_count gauge
prometheus_toolbox_remote_write_timeseries_count{origin="[::1]:58060"} 2
# HELP prometheus_toolbox_remote_write_uncompressed_size uncompressed remote write request size in bytes
# TYPE prometheus_toolbox_remote_write_uncompressed_size gauge
prometheus_toolbox_remote_write_uncompressed_size{origin="[::1]:58060"} 143
```
