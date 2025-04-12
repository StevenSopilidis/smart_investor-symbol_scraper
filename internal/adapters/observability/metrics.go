package observability

import (
	"context"
	"time"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitMeterProvider(ctx context.Context, config config.Config) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.OtelCollectorAddr),
	)
	if err != nil {
		return nil, err
	}

	pariodicReader := metric.NewPeriodicReader(exp, metric.WithInterval(10*time.Second))

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(pariodicReader),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}
