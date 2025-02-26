# Benchmarks for OTel Go Span Events vs Logs

## Usage

```sh
make up
make stat
make down
```

`make up` runs the OpenTelemetry Collector container in background. Give it some time.

`make bench` runs the Go benchmarks.

`make stat` runs the Go benchmarks and computes their statistical summaries and
A/B comparisons using [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat).

`make down` stops and removes the OpenTelemetry Collector container.

You can also run everything:

```sh
make
```

Example output:

```out
docker start otelcol || docker run -d --name otelcol -p 127.0.0.1:4318:4318 -p 127.0.0.1:55679:55679 otel/opentelemetry-collector:0.120.0
otelcol
go test -run=^$ -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/pellared/spanevents-vs-logs
cpu: 13th Gen Intel(R) Core(TM) i7-13800H
BenchmarkSpanEvents/OTLP-20                24157             43275 ns/op          103485 B/op        739 allocs/op
BenchmarkSpanEvents/STDOUT-20              29248             38446 ns/op           90638 B/op        567 allocs/op
BenchmarkSpanEvents/NOOP-20                74344             15907 ns/op           29840 B/op        303 allocs/op
BenchmarkLogs/OTLP-20                       3486           3969658 ns/op          189184 B/op       1770 allocs/op
BenchmarkLogs/STDOUT-20                     1257            970677 ns/op          281219 B/op       4608 allocs/op
BenchmarkLogs/NOOP-20                     207175              5703 ns/op             240 B/op          3 allocs/op
PASS
ok      github.com/pellared/spanevents-vs-logs  20.482s
go test -run=^$ -bench=BenchmarkSpanEvents -benchmem -count=10 > spanevents.out 
sed -i -e 's/BenchmarkSpanEvents/Benchmark/g' spanevents.out
go test -run=^$ -bench=BenchmarkLogs -benchmem -count=10 > logs.out
sed -i -e 's/BenchmarkLogs/Benchmark/g' logs.out
go tool benchstat spanevents.out logs.out
goos: linux
goarch: amd64
pkg: github.com/pellared/spanevents-vs-logs
cpu: 13th Gen Intel(R) Core(TM) i7-13800H
           │ spanevents.out │                 logs.out                 │
           │     sec/op     │     sec/op      vs base                  │
/OTLP-20       41.27µ ±  4%   1861.14µ ± 11%  +4409.24% (p=0.000 n=10)
/STDOUT-20     38.04µ ±  2%    932.10µ ±  4%  +2350.38% (p=0.000 n=10)
/NOOP-20      27.532µ ± 39%     5.727µ ±  1%    -79.20% (p=0.000 n=10)
geomean        35.09µ           215.0µ         +512.55%

           │ spanevents.out │                logs.out                │
           │      B/op      │     B/op       vs base                 │
/OTLP-20       100.9Ki ± 2%    196.2Ki ± 0%   +94.44% (p=0.000 n=10)
/STDOUT-20     88.59Ki ± 0%   273.90Ki ± 1%  +209.16% (p=0.000 n=10)
/NOOP-20       29840.0 ± 0%      240.0 ± 0%   -99.20% (p=0.000 n=10)
geomean        63.86Ki         23.26Ki        -63.57%

           │ spanevents.out │               logs.out               │
           │   allocs/op    │  allocs/op   vs base                 │
/OTLP-20         737.0 ± 6%   1905.0 ± 0%  +158.48% (p=0.000 n=10)
/STDOUT-20       569.0 ± 1%   4603.5 ± 1%  +709.05% (p=0.000 n=10)
/NOOP-20       303.000 ± 0%    3.000 ± 0%   -99.01% (p=0.000 n=10)
geomean          502.7         297.4        -40.84%
```
 