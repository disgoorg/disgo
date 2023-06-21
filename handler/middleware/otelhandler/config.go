package otelhandler

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func DefaultConfig() *Config {
	return &Config{
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
		Filter:         nil,
	}
}

type Config struct {
	TracerProvider oteltrace.TracerProvider
	MeterProvider  otelmetric.MeterProvider
	Filter         func(ctx context.Context, e *events.InteractionCreate) bool
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithTracerProvider(provider oteltrace.TracerProvider) ConfigOpt {
	return func(cfg *Config) {
		cfg.TracerProvider = provider
	}
}

func WithMeterProvider(provider otelmetric.MeterProvider) ConfigOpt {
	return func(cfg *Config) {
		cfg.MeterProvider = provider
	}
}

func WithFilter(filter func(ctx context.Context, e *events.InteractionCreate) bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.Filter = filter
	}
}
