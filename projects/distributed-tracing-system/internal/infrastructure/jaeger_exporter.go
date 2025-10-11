package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

// jaegerExporter implements the JaegerExporter interface
type jaegerExporter struct {
	exporter *jaeger.Exporter
	client   *http.Client
	endpoint string
}

// NewJaegerExporter creates a new Jaeger exporter
func NewJaegerExporter(endpoint string) (domain.JaegerExporter, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("jaeger endpoint is required")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(endpoint),
		jaeger.WithHTTPClient(client),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	return &jaegerExporter{
		exporter: exporter,
		client:   client,
		endpoint: endpoint,
	}, nil
}

// ExportTrace exports a trace to Jaeger
func (je *jaegerExporter) ExportTrace(ctx context.Context, trace *domain.Trace) error {
	// Validate trace
	if err := je.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
	}

	// Convert domain trace to OpenTelemetry span
	spanData, err := je.convertTraceToSpanData(trace)
	if err != nil {
		return fmt.Errorf("failed to convert trace to span data: %w", err)
	}

	// Export to Jaeger
	if err := je.exporter.ExportSpans(ctx, []trace.ReadOnlySpan{spanData}); err != nil {
		return fmt.Errorf("failed to export trace to Jaeger: %w", err)
	}

	return nil
}

// validateTrace validates a trace before export
func (je *jaegerExporter) validateTrace(trace *domain.Trace) error {
	if trace.ID == "" {
		return fmt.Errorf("trace ID is required")
	}
	if trace.Service == "" {
		return fmt.Errorf("service name is required")
	}
	if trace.Operation == "" {
		return fmt.Errorf("operation name is required")
	}
	if trace.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if trace.EndTime.IsZero() {
		return fmt.Errorf("end time is required")
	}
	if trace.StartTime.After(trace.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}
	return nil
}

// convertTraceToSpanData converts domain trace to OpenTelemetry span data
func (je *jaegerExporter) convertTraceToSpanData(trace *domain.Trace) (trace.ReadOnlySpan, error) {
	// This is a simplified conversion
	// In a real implementation, you would properly convert the domain trace
	// to OpenTelemetry span data structure
	
	// For now, return a mock span data
	// This would need proper implementation with OpenTelemetry SDK
	return nil, fmt.Errorf("conversion not implemented yet")
}

