package metrics

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitMetrics() *metric.MeterProvider {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.DeploymentEnvironmentKey.String("development"),
		),
	)
	if err != nil {
		log.Printf("warning: failed to create otel resource: %v", err)
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		log.Printf("warning: failed to create metric exporter: %v", err)
		return nil
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(15*time.Second),
			),
		),
	)
	otel.SetMeterProvider(provider)

	err = runtime.Start(runtime.WithMeterProvider(provider))
	if err != nil {
		log.Printf("warning: failed to start runtime metrics: %v", err)
	}

	return provider
}
