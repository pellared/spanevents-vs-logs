package bench

import (
	"context"
	"io"
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
	traceapi "go.opentelemetry.io/otel/trace"
)

func BenchmarkSpanEvents(b *testing.B) {
	testCases := []struct {
		name       string
		traceExpFn func(b *testing.B) trace.SpanExporter
	}{
		{
			name:       "OTLP",
			traceExpFn: setupOTLPTraceExporter,
		},
		{
			name:       "STDOUT",
			traceExpFn: setupSTDOUTTraceExporter,
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracerProvider := trace.NewTracerProvider(
				trace.WithBatcher(tc.traceExpFn(b)))
			b.Cleanup(func() {
				if err := tracerProvider.Shutdown(context.Background()); err != nil {
					b.Fatalf("tracerProvider.Shutdown: %v", err)
				}
			})
			tracer := tracerProvider.Tracer(b.Name())

			for b.Loop() {
				// code to measure
				_, span := tracer.Start(context.Background(), b.Name())
				for i := range 100 {
					msg := strconv.Itoa(i)
					span.AddEvent(msg, traceapi.WithAttributes(
						attribute.Bool("b", true),
						attribute.Float64("pi", 3.14),
						attribute.Int("ten", 10),
						attribute.String("foo", "bar"),
					))
				}
				span.End()
			}
		})
	}

}

func BenchmarkLogs(b *testing.B) {
	testCases := []struct {
		name       string
		traceExpFn func(b *testing.B) trace.SpanExporter
		logExpFn   func(b *testing.B) log.Exporter
	}{
		{
			name:       "OTLP",
			traceExpFn: setupOTLPTraceExporter,
			logExpFn:   setupOTLPLogExporter,
		},
		{
			name:       "STDOUT",
			traceExpFn: setupSTDOUTTraceExporter,
			logExpFn:   setupSTDOUTLogExporter,
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracerProvider := trace.NewTracerProvider(
				trace.WithBatcher(tc.traceExpFn(b)))
			b.Cleanup(func() {
				if err := tracerProvider.Shutdown(context.Background()); err != nil {
					b.Fatalf("tracerProvider.Shutdown: %v", err)
				}
			})
			tracer := tracerProvider.Tracer(b.Name())

			logProvider := log.NewLoggerProvider(
				log.WithProcessor(log.NewBatchProcessor(tc.logExpFn(b))))
			b.Cleanup(func() {
				if err := logProvider.Shutdown(context.Background()); err != nil {
					b.Fatalf("logProvider.Shutdown: %v", err)
				}
			})
			logger := logProvider.Logger(b.Name())

			for b.Loop() {
				// code to measure
				ctx, span := tracer.Start(b.Context(), b.Name())
				for i := range 100 {
					msg := strconv.Itoa(i)
					var r logapi.Record
					r.SetBody(logapi.StringValue(msg))
					r.AddAttributes(
						logapi.Bool("b", true),
						logapi.Float64("pi", 3.14),
						logapi.Int("ten", 10),
						logapi.String("foo", "bar"),
					)
					logger.Emit(ctx, r)
				}
				span.End()
			}
		})
	}
}

func setupOTLPTraceExporter(b *testing.B) trace.SpanExporter {
	exp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlptracehttp.New: %v", err)
	}
	return exp
}

func setupSTDOUTTraceExporter(b *testing.B) trace.SpanExporter {
	exp, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	if err != nil {
		b.Fatalf("stdouttrace.New: %v", err)
	}
	return exp
}

func setupOTLPLogExporter(b *testing.B) log.Exporter {
	exp, err := otlploghttp.New(b.Context(), otlploghttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlploghttp.New: %v", err)
	}
	return exp
}

func setupSTDOUTLogExporter(b *testing.B) log.Exporter {
	exp, err := stdoutlog.New(stdoutlog.WithWriter(io.Discard))
	if err != nil {
		b.Fatalf("stdoutlog.New: %v", err)
	}
	return exp
}
