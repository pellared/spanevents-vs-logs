package bench

import (
	"context"
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
	traceapi "go.opentelemetry.io/otel/trace"
)

func BenchmarkSpanEvents(b *testing.B) {
	traceExporter, err := otlptracehttp.New(b.Context(), otlptracehttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlptracehttp.New: %v", err)
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter))
	b.Cleanup(func() {
		if err := traceExporter.Shutdown(context.Background()); err != nil {
			b.Fatalf("traceExporter.Shutdown: %v", err)
		}
	})
	tracer := tracerProvider.Tracer(b.Name())

	for b.Loop() {
		for range 100 {
			_, span := tracer.Start(b.Context(), b.Name())
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

		if err := tracerProvider.ForceFlush(b.Context()); err != nil {
			b.Fatalf("tracerProvider.ForceFlush: %v", err)
		}
	}
}

func BenchmarkLogs(b *testing.B) {
	traceExporter, err := otlptracehttp.New(b.Context(), otlptracehttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlptracehttp.New: %v", err)
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter))
	b.Cleanup(func() {
		if err := traceExporter.Shutdown(context.Background()); err != nil {
			b.Fatalf("traceExporter.Shutdown: %v", err)
		}
	})
	tracer := tracerProvider.Tracer(b.Name())

	logExporter, err := otlploghttp.New(b.Context(), otlploghttp.WithInsecure())
	if err != nil {
		b.Fatalf("otlploghttp.New: %v", err)
	}
	logProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)))
	b.Cleanup(func() {
		if err := logProvider.Shutdown(context.Background()); err != nil {
			b.Fatalf("logProvider.Shutdown: %v", err)
		}
	})
	logger := logProvider.Logger(b.Name())

	for b.Loop() {
		for range 100 {
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

		if err := tracerProvider.ForceFlush(b.Context()); err != nil {
			b.Fatalf("tracerProvider.ForceFlush: %v", err)
		}
		if err := logProvider.ForceFlush(b.Context()); err != nil {
			b.Fatalf("logProvider.ForceFlush: %v", err)
		}
	}
}
