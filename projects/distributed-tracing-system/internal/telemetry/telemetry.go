package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryConfig holds configuration for telemetry
type TelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
	PrometheusPort string
	SamplingRate   float64
}

// TelemetryManager manages OpenTelemetry SDK components
type TelemetryManager struct {
	config        *TelemetryConfig
	tracer        trace.Tracer
	meter         metric.Meter
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
}

// NewTelemetryManager creates a new telemetry manager
func NewTelemetryManager(config *TelemetryConfig) (*TelemetryManager, error) {
	if config == nil {
		return nil, fmt.Errorf("telemetry config is required")
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", config.ServiceName),
			attribute.String("service.version", config.ServiceVersion),
			attribute.String("deployment.environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create Jaeger exporter
	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(config.JaegerEndpoint),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create Prometheus exporter
	prometheusExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SamplingRate)),
	)

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(prometheusExporter),
		sdkmetric.WithResource(res),
	)

	// Set global providers
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer and meter
	tracer := tracerProvider.Tracer(config.ServiceName)
	meter := meterProvider.Meter(config.ServiceName)

	return &TelemetryManager{
		config:         config,
		tracer:         tracer,
		meter:          meter,
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
	}, nil
}

// GetTracer returns the tracer
func (tm *TelemetryManager) GetTracer() trace.Tracer {
	return tm.tracer
}

// GetMeter returns the meter
func (tm *TelemetryManager) GetMeter() metric.Meter {
	return tm.meter
}

// StartSpan starts a new span
func (tm *TelemetryManager) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tm.tracer.Start(ctx, name, opts...)
}

// StartSpanWithAttributes starts a new span with attributes
func (tm *TelemetryManager) StartSpanWithAttributes(ctx context.Context, name string, attrs []attribute.KeyValue, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	opts = append(opts, trace.WithAttributes(attrs...))
	return tm.tracer.Start(ctx, name, opts...)
}

// RecordMetric records a metric
func (tm *TelemetryManager) RecordMetric(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue) error {
	// Create a counter metric
	counter, err := tm.meter.Float64Counter(name)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}

	counter.Add(ctx, value, metric.WithAttributes(attrs...))
	return nil
}

// RecordHistogram records a histogram metric
func (tm *TelemetryManager) RecordHistogram(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue) error {
	// Create a histogram metric
	histogram, err := tm.meter.Float64Histogram(name)
	if err != nil {
		return fmt.Errorf("failed to create histogram: %w", err)
	}

	histogram.Record(ctx, value, metric.WithAttributes(attrs...))
	return nil
}

// Shutdown gracefully shuts down the telemetry manager
func (tm *TelemetryManager) Shutdown(ctx context.Context) error {
	var err error

	// Shutdown tracer provider
	if tm.tracerProvider != nil {
		if shutdownErr := tm.tracerProvider.Shutdown(ctx); shutdownErr != nil {
			err = fmt.Errorf("failed to shutdown tracer provider: %w", shutdownErr)
		}
	}

	// Shutdown meter provider
	if tm.meterProvider != nil {
		if shutdownErr := tm.meterProvider.Shutdown(ctx); shutdownErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to shutdown meter provider: %w; previous error: %v", shutdownErr, err)
			} else {
				err = fmt.Errorf("failed to shutdown meter provider: %w", shutdownErr)
			}
		}
	}

	return err
}

// GetDefaultConfig returns a default telemetry configuration
func GetDefaultConfig() *TelemetryConfig {
	return &TelemetryConfig{
		ServiceName:    "distributed-tracing-system",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		JaegerEndpoint: "http://localhost:14268/api/traces",
		PrometheusPort: "9091",
		SamplingRate:   0.1, // 10% sampling
	}
}
