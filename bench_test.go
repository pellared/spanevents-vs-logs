package bench

import (
	"context"
	"io"
	"os"
	"path/filepath"
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
		name  string
		expFn func(b *testing.B) trace.SpanExporter
	}{
		{
			name: "OTLP",
			expFn: func(b *testing.B) trace.SpanExporter {
				traceExporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())
				if err != nil {
					b.Fatalf("otlptracehttp.New: %v", err)
				}
				return traceExporter
			},
		},
		{
			name: "STDOUT",
			expFn: func(b *testing.B) trace.SpanExporter {
				traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
				if err != nil {
					b.Fatalf("stdouttrace.New: %v", err)
				}
				return traceExporter
			},
		},
		{
			name: "FILE",
			expFn: func(b *testing.B) trace.SpanExporter {
				f, err := os.Create(filepath.Join(b.TempDir(), "traces.out"))
				if err != nil {
					b.Fatalf("os.Create: %v", err)
				}
				b.Cleanup(func() {
					f.Close()
				})
				traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(f))
				if err != nil {
					b.Fatalf("stdouttrace.New: %v", err)
				}
				return traceExporter
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracerProvider := trace.NewTracerProvider(
				trace.WithBatcher(tc.expFn(b)))
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
		name  string
		expFn func(b *testing.B) (trace.SpanExporter, log.Exporter)
	}{
		{
			name: "OTLP",
			expFn: func(b *testing.B) (trace.SpanExporter, log.Exporter) {
				traceExporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())
				if err != nil {
					b.Fatalf("otlptracehttp.New: %v", err)
				}
				logExporter, err := otlploghttp.New(b.Context(), otlploghttp.WithInsecure())
				if err != nil {
					b.Fatalf("otlploghttp.New: %v", err)
				}
				return traceExporter, logExporter
			},
		},
		{
			name: "STDOUT",
			expFn: func(b *testing.B) (trace.SpanExporter, log.Exporter) {
				traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
				if err != nil {
					b.Fatalf("stdouttrace.New: %v", err)
				}
				logExporter, err := stdoutlog.New(stdoutlog.WithWriter(io.Discard))
				if err != nil {
					b.Fatalf("stdoutlog.New: %v", err)
				}
				return traceExporter, logExporter
			},
		},
		{
			name: "FILE",
			expFn: func(b *testing.B) (trace.SpanExporter, log.Exporter) {
				f, err := os.Create(filepath.Join(b.TempDir(), "traces.out"))
				if err != nil {
					b.Fatalf("os.Create: %v", err)
				}
				b.Cleanup(func() {
					f.Close()
				})
				traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(f))
				if err != nil {
					b.Fatalf("stdouttrace.New: %v", err)
				}
				f, err = os.Create(filepath.Join(b.TempDir(), "logs.out"))
				if err != nil {
					b.Fatalf("os.Create: %v", err)
				}
				b.Cleanup(func() {
					f.Close()
				})
				logExporter, err := stdoutlog.New(stdoutlog.WithWriter(f))
				if err != nil {
					b.Fatalf("stdoutlog.New: %v", err)
				}
				return traceExporter, logExporter
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tExp, lExp := tc.expFn(b)

			tracerProvider := trace.NewTracerProvider(
				trace.WithBatcher(tExp))
			b.Cleanup(func() {
				if err := tracerProvider.Shutdown(context.Background()); err != nil {
					b.Fatalf("tracerProvider.Shutdown: %v", err)
				}
			})
			tracer := tracerProvider.Tracer(b.Name())

			logProvider := log.NewLoggerProvider(
				log.WithProcessor(log.NewBatchProcessor(lExp)))
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
