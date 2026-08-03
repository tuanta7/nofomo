package o11y

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

type MeterProvider struct {
	service  string
	provider *sdkmetric.MeterProvider
}

func NewMeterProvider(ctx context.Context, service string) (*MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(grpcExportEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(service)),
	)
	if err != nil {
		return nil, err
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))
	return &MeterProvider{
		service: service,
		provider: sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(reader),
		),
	}, nil
}

func (p *MeterProvider) Meter(serviceName ...string) metric.Meter {
	serviceName = append(serviceName, p.service)
	return p.provider.Meter(serviceName[0])
}

func (p *MeterProvider) Shutdown(ctx context.Context) error {
	if p.provider == nil {
		return nil
	}

	return p.provider.Shutdown(ctx)
}
