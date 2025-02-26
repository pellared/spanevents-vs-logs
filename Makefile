bench:
	go test -run=^$$ -bench=. -benchmem

stat:
	go test -run=^$$ -bench=BenchmarkSpanEvents -benchmem -count=10 > spanevents.out 
	sed -i -e 's/BenchmarkSpanEvents/Benchmark/g' spanevents.out
	go test -run=^$$ -bench=BenchmarkLogs -benchmem -count=10 > logs.out
	sed -i -e 's/BenchmarkLogs/Benchmark/g' logs.out
	go tool benchstat spanevents.out logs.out
