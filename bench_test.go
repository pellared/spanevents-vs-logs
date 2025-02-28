package bench

import (
	"context"
	"io"
	"strconv"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
	traceapi "go.opentelemetry.io/otel/trace"
)

func BenchmarkSpan(b *testing.B) {
	testCases := []struct {
		name     string
		tracerFn func(b *testing.B) traceapi.Tracer
	}{
		{
			name: "OTLP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupOTLPTraceExporter(b)
				return setupTracing(b, traceExp)
			},
		},
		{
			name: "STDOUT",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupSTDOUTTraceExporter(b)
				return setupTracing(b, traceExp)
			},
		},
		{
			name: "NOOP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				return otel.Tracer(b.Name())
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracer := tc.tracerFn(b)

			for b.Loop() {
				// code to measure
				_, span := tracer.Start(context.Background(), b.Name())
				span.End()
			}
		})
	}

}

func BenchmarkLog(b *testing.B) {
	testCases := []struct {
		name     string
		loggerFn func(b *testing.B) logapi.Logger
	}{
		{
			name: "OTLP",
			loggerFn: func(b *testing.B) logapi.Logger {
				logExp := setupOTLPLogExporter(b)
				return setupLogging(b, logExp)
			},
		},
		{
			name: "STDOUT",
			loggerFn: func(b *testing.B) logapi.Logger {
				logExp := setupSTDOUTLogExporter(b)
				return setupLogging(b, logExp)
			},
		},
		{
			name: "NOOP",
			loggerFn: func(b *testing.B) logapi.Logger {
				return global.Logger(b.Name())
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			logger := tc.loggerFn(b)

			for b.Loop() {
				// code to measure
				var r logapi.Record
				logger.Emit(context.Background(), r)
			}
		})
	}
}

func BenchmarkSpanEvents(b *testing.B) {
	testCases := []struct {
		name     string
		tracerFn func(b *testing.B) traceapi.Tracer
	}{
		{
			name: "OTLP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupOTLPTraceExporter(b)
				return setupTracing(b, traceExp)
			},
		},
		{
			name: "STDOUT",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupSTDOUTTraceExporter(b)
				return setupTracing(b, traceExp)
			},
		},
		{
			name: "NOOP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				return otel.Tracer(b.Name())
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracer := tc.tracerFn(b)

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
		name     string
		tracerFn func(b *testing.B) traceapi.Tracer
		loggerFn func(b *testing.B) logapi.Logger
	}{
		{
			name: "OTLP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupOTLPTraceExporter(b)
				return setupTracing(b, traceExp)
			},
			loggerFn: func(b *testing.B) logapi.Logger {
				logExp := setupOTLPLogExporter(b)
				return setupLogging(b, logExp)
			},
		},
		{
			name: "STDOUT",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				traceExp := setupSTDOUTTraceExporter(b)
				return setupTracing(b, traceExp)
			},
			loggerFn: func(b *testing.B) logapi.Logger {
				logExp := setupSTDOUTLogExporter(b)
				return setupLogging(b, logExp)
			},
		},
		{
			name: "NOOP",
			tracerFn: func(b *testing.B) traceapi.Tracer {
				return otel.Tracer(b.Name())
			},
			loggerFn: func(b *testing.B) logapi.Logger {
				return global.Logger(b.Name())
			},
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// setup
			tracer := tc.tracerFn(b)
			logger := tc.loggerFn(b)

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

func setupTracing(b *testing.B, exp trace.SpanExporter) traceapi.Tracer {
	b.Helper()
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exp))
	b.Cleanup(func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			b.Fatalf("tracerProvider.Shutdown: %v", err)
		}
	})
	return tracerProvider.Tracer(b.Name())
}

func setupLogging(b *testing.B, exp log.Exporter) logapi.Logger {
	b.Helper()
	logProvider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(exp)))
	b.Cleanup(func() {
		if err := logProvider.Shutdown(context.Background()); err != nil {
			b.Fatalf("logProvider.Shutdown: %v", err)
		}
	})
	return logProvider.Logger(b.Name())
}

func setupOTLPTraceExporter(b *testing.B) trace.SpanExporter {
	b.Helper()
	exp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlptracehttp.New: %v", err)
	}
	return exp
}

func setupSTDOUTTraceExporter(b *testing.B) trace.SpanExporter {
	b.Helper()
	exp, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	if err != nil {
		b.Fatalf("stdouttrace.New: %v", err)
	}
	return exp
}

func setupOTLPLogExporter(b *testing.B) log.Exporter {
	b.Helper()
	exp, err := otlploghttp.New(b.Context(), otlploghttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlploghttp.New: %v", err)
	}
	return exp
}

func setupSTDOUTLogExporter(b *testing.B) log.Exporter {
	b.Helper()
	exp, err := stdoutlog.New(stdoutlog.WithWriter(io.Discard))
	if err != nil {
		b.Fatalf("stdoutlog.New: %v", err)
	}
	return exp
}
