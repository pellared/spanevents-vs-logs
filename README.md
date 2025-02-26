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

 ```
 docker start otelcol || docker run -d --name otelcol -p 127.0.0.1:4318:4318 -p 127.0.0.1:55679:55679 otel/opentelemetry-collector:0.120.0
otelcol
go test -run=^$ -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/pellared/spanevents-vs-logs
cpu: 13th Gen Intel(R) Core(TM) i7-13800H
BenchmarkSpanEvents/OTLP-20                22708             45040 ns/op          104856 B/op        765 allocs/op
BenchmarkSpanEvents/STDOUT-20              28063             40038 ns/op           90748 B/op        569 allocs/op
BenchmarkLogs/OTLP-20                       3309           4516423 ns/op          184387 B/op       1717 allocs/op
BenchmarkLogs/STDOUT-20                     1396            975572 ns/op          279547 B/op       4607 allocs/op
PASS
ok      github.com/pellared/spanevents-vs-logs  19.371s
go test -run=^$ -bench=BenchmarkSpanEvents -benchmem -count=10 > spanevents.out 
sed -i -e 's/BenchmarkSpanEvents/Benchmark/g' spanevents.out
go test -run=^$ -bench=BenchmarkLogs -benchmem -count=10 > logs.out
sed -i -e 's/BenchmarkLogs/Benchmark/g' logs.out
go tool benchstat spanevents.out logs.out
goos: linux
goarch: amd64
pkg: github.com/pellared/spanevents-vs-logs
cpu: 13th Gen Intel(R) Core(TM) i7-13800H
           │ spanevents.out │                logs.out                 │
           │     sec/op     │    sec/op      vs base                  │
/OTLP-20        41.72µ ± 4%   2023.12µ ± 9%  +4749.04% (p=0.000 n=10)
/STDOUT-20      38.72µ ± 1%    945.33µ ± 4%  +2341.54% (p=0.000 n=10)
geomean         40.19µ          1.383m       +3340.81%

           │ spanevents.out │                logs.out                │
           │      B/op      │     B/op       vs base                 │
/OTLP-20       101.0Ki ± 2%    195.8Ki ± 1%   +93.81% (p=0.000 n=10)
/STDOUT-20     88.71Ki ± 0%   272.19Ki ± 0%  +206.81% (p=0.000 n=10)
geomean        94.67Ki         230.9Ki       +143.85%

           │ spanevents.out │               logs.out               │
           │   allocs/op    │  allocs/op   vs base                 │
/OTLP-20         739.0 ± 5%   1901.5 ± 1%  +157.31% (p=0.000 n=10)
/STDOUT-20       572.5 ± 0%   4598.5 ± 0%  +703.23% (p=0.000 n=10)
geomean          650.4        2.957k       +354.62%
 ```
 