# Prometheus Toolbox

A CLI tool to generate metrics data for the [Prometheus](https://prometheus.io/) monitorig system.
It allows you to define time series and values and send them directly to Prometheus using the remote write protocol.

## Why?

I found this very useful for testing PromQL queries and alerts.
Sometimes you don't have the right data in Prometheus or simply no access. 
This tool let's you simulate the data you need and test your queries against it.

## Installation

The latest version:

`GOPROXY=direct go install github.com/pb82/prometheus-toolbox@latest`

A specific version:

`GOPROXY=direct go install github.com/pb82/prometheus-toolbox@<version>`

## Starting a local Prometheus instance

You can use the included Makefile target:

```shell
$ make prom
```

or directly with docker:

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 10s
  evaluation_interval: 30s
```

then run the container:

`docker run -p 9090:9090 --network=host -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus --web.enable-remote-write-receiver --config.file=/etc/prometheus/prometheus.yml`

# Flags

The following flags are accepted:

* `--help` Print help and exit
* `--version` Print version and exit
* `--prometheus.url` Prometheus base url
* `--config.file` Config file location
* `--batch.size` Samples per remote write request, defaults to 500

# Config file format

This tool reads from a config file where the simulated time series and values are defined.
The format is:

```yaml
interval: "10s"                   # Interval between samples, in this case 10 seconds
time_series:                      # List of time series to simulate
  - series: metric_a{label="a"}   # Time series (metric name and label list)
    values: 1+1x100               # Precalculated samples
  - series: metric_a{label="a"}   # Another time series
    stream: 1+0                   # Realtime samples
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

# Features

[Data generator](#data-generator): write precalculated samples.

[Streaming samples](#streaming-samples): write realtime data.

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
