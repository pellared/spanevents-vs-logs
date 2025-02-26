all: bench stat

stat: up
	go test -run=^$$ -bench=BenchmarkSpanEvents -benchmem -count=10 > spanevents.out 
	sed -i -e 's/BenchmarkSpanEvents/Benchmark/g' spanevents.out
	go test -run=^$$ -bench=BenchmarkLogs -benchmem -count=10 > logs.out
	sed -i -e 's/BenchmarkLogs/Benchmark/g' logs.out
	go tool benchstat spanevents.out logs.out

bench: up
	go test -run=^$$ -bench=. -benchmem

up:
	docker start otelcol || docker run -d --name otelcol -p 127.0.0.1:4318:4318 -p 127.0.0.1:55679:55679 otel/opentelemetry-collector:0.120.0

down:
	docker stop otelcol
	docker rm otelcol
