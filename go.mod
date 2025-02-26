module github.com/pellared/spanevents-vs-logs

go 1.24.0

tool golang.org/x/perf/cmd/benchstat

require (
	go.opentelemetry.io/otel v1.34.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.10.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.34.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.10.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.34.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/log v0.10.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/sdk v1.34.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/sdk/log v0.10.1-0.20250226060240-86d783c3daf7
	go.opentelemetry.io/otel/trace v1.34.1-0.20250226060240-86d783c3daf7
)

require (
	github.com/aclements/go-moremath v0.0.0-20210112150236-f10218a38794 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.1-0.20250226060240-86d783c3daf7 // indirect
	go.opentelemetry.io/otel/metric v1.34.1-0.20250226060240-86d783c3daf7 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/perf v0.0.0-20250214215153-c95ad7d5b636 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.70.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
