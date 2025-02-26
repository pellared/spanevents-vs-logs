# Benchmarks for OTel Go Span Events vs Logs

## Usage

```sh
make up
make bench
make down
```

`make up` runs the OpenTelemetry Collector container in background. Give it some time.

`make stat` runs the Go benchmarks and computes their statistical summaries and
A/B comparisons using [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat).

`make down` stops and removes the OpenTelemetry Collector container.
 