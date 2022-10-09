package metrics

import (
	"context"
	"fmt"

	"github.com/caarlos0/env"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer("mailsender")

type MetricsContainer struct {
	*jaeger.Exporter
	*tracesdk.TracerProvider
}

type Config struct {
	ServiceName   string `env:"SERVICE_NAME"`
	JaegerAddress string `env:"JAEGER_ADDRESS" envDefault:"http://jaeger:14268/api/traces"`
}

func New() (*MetricsContainer, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse metrics configuration failed: %w", err)
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(cfg.JaegerAddress)),
	)
	if err != nil {
		return nil, fmt.Errorf("init jaeger failed: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &MetricsContainer{
		Exporter:       exp,
		TracerProvider: tp,
	}, nil

}

func (mc *MetricsContainer) Stop(ctx context.Context) error {
	err := mc.Exporter.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("jaeger shutdown failed: %w", err)
	}

	err = mc.TracerProvider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("tracer proviser shutdown failed: %w", err)
	}

	return nil
}
