# Prometheus Toolbox

A CLI tool to interact with the [Prometheus](https://prometheus.io/) monitorig system.
Use the included Makefile to start a local instance on port `9090`:

```shell
$ make prom
```

# Features

This is a work in progress, additional features for testing and debugging Queries and Alerts will be added.

## Data generator

Generates data and sends it to a Prometheus instance using the [remote write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) protocol.

**NOTE:** remote write needs to be specifically enabled in Prometheus with the flag `--web.enable-remote-write-receiver`

Provide the base URL of the Prometheus instance and a config file, for example:

```shell
`$ cat example/config.yaml`
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

`$ ./prometheus-toolbox --prometheus.url=http://localhost:9090 --config.file=./example/config.yml`

Open the Prometheus UI and run a query, e.g. `{__name__~="metric_a|metric_b""}` or go to the graph view.
Samples for both time series should be visible.

### How the samples are generated

Samples are always generated for the past. The time series that produces the most samples (in the above example it is `metrics_a`) determines how far back in time the samples go.
In the above case 100 samples are generated with an interval of ten seconds between them.
So both series will start at a point in time 1000 seconds in the past.
