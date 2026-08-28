package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go-backend-boilerplate/internal/config"
)

type ShutdownFunc func(context.Context) error

func Setup(ctx context.Context, cfg *config.Config) (ShutdownFunc, error) {
	noop := func(context.Context) error { return nil }
	if !cfg.OTelEnabled {
		return noop, nil
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.OTelEndpoint)}
	if cfg.OTelInsecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return noop, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", cfg.OTelServiceName),
	))
	if err != nil {
		return noop, err
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
