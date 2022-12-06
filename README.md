# Prometheus Toolbox

A CLI tool to interact with the [Prometheus](https://prometheus.io/) monitorig system.
Use the included Makefile to start a local instance on port `9090`:

```shell
$ make prom
```

To build the CLI, run:

```shell
$ make build
```

# Flags

The following flags are accepted:

* `--version` Print version and exit
* `--prometheus.url` Prometheus base url
* `--config.file` Config file location
* `--batch.size` Samples per remote write request, defaults to 500

# Features

This is a work in progress, additional features for testing and debugging Queries and Alerts will be added.

## Data generator

Generates data and sends it to a Prometheus instance using the [remote write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) protocol.

**NOTE:** remote write needs to be specifically enabled in Prometheus with the flag `--web.enable-remote-write-receiver`

Provide the base URL of the Prometheus instance and a config file, for example:

```shell
$ cat example/config.yaml
---
interval: 10s
time_series:
  - series: metric_a{label="a"}
    values: 1+0x100
  - series: metric_b{label="b"}
    values: 1+1x50
```

**NOTE:** if you're not familiar with the `series` and `values` notation, it is similar to how unit tests are written. For more information, have a look [here](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/#series).

Run the CLI with the following options:

```shell
$ ./promtoolbox --prometheus.url=http://localhost:9090 --config.file=./example/config.yml
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