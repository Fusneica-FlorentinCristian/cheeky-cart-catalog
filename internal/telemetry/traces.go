package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const defaultServiceName = "cheeky-cart-catalog-rest"

// InitTracing configures OTLP HTTP export when OTEL_EXPORTER_OTLP_ENDPOINT is set.
// Returns a shutdown func (no-op when tracing is disabled).
func InitTracing(ctx context.Context) (func(context.Context) error, error) {
	endpoint := strings.TrimPrefix(strings.TrimPrefix(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "http://"), "https://")
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	protocol := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	var exporter sdktrace.SpanExporter
	var err error
	if protocol == "grpc" {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
	} else {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("otlp exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// TraceMiddleware wraps h with otelhttp when OTEL_EXPORTER_OTLP_ENDPOINT is set.
func TraceMiddleware(h http.Handler) http.Handler {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return h
	}
	return otelhttp.NewHandler(h, "catalog-rest",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// FlushTracing gives the batch exporter time to send spans before exit.
func FlushTracing(ctx context.Context, shutdown func(context.Context) error) {
	if shutdown == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = shutdown(ctx)
}
