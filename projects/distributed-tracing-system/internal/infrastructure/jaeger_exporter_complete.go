package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/streamforge/distributed-tracing-system/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// jaegerExporterComplete implements the JaegerExporter interface using OpenTelemetry SDK
type jaegerExporterComplete struct {
	telemetryManager *telemetry.TelemetryManager
	tracer          oteltrace.Tracer
}

// NewJaegerExporterComplete creates a new complete Jaeger exporter
func NewJaegerExporterComplete(telemetryManager *telemetry.TelemetryManager) (domain.JaegerExporter, error) {
	if telemetryManager == nil {
		return nil, fmt.Errorf("telemetry manager is required")
	}

	return &jaegerExporterComplete{
		telemetryManager: telemetryManager,
		tracer:          telemetryManager.GetTracer(),
	}, nil
}

// ExportTrace exports a trace to Jaeger using OpenTelemetry SDK
func (je *jaegerExporterComplete) ExportTrace(ctx context.Context, trace *domain.Trace) error {
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

	// Create a span for the trace export
	spanName := fmt.Sprintf("export-trace-%s", trace.Operation)
	ctx, span := je.tracer.Start(ctx, spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		oteltrace.WithAttributes(
			attribute.String("trace.id", string(trace.ID)),
			attribute.String("trace.service", string(trace.Service)),
			attribute.String("trace.operation", string(trace.Operation)),
			attribute.String("trace.status", string(trace.Status)),
			attribute.Int64("trace.duration_ns", trace.Duration.Nanoseconds()),
			attribute.Int("trace.spans_count", len(trace.Spans)),
		),
	)
	defer span.End()

	// Add custom tags as attributes
	for key, value := range trace.Tags {
		span.SetAttributes(attribute.String("trace.tag."+key, value))
	}

	// Process spans
	if err := je.processSpans(ctx, span, trace); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to process spans: %w", err)
	}

	// Record metrics
	je.recordTraceMetrics(ctx, trace)

	// Set span status based on trace status
	switch trace.Status {
	case domain.TraceStatusSuccess:
		span.SetStatus(codes.Ok, "Trace exported successfully")
	case domain.TraceStatusError:
		span.SetStatus(codes.Error, "Trace contains errors")
	case domain.TraceStatusTimeout:
		span.SetStatus(codes.Error, "Trace timed out")
	default:
		span.SetStatus(codes.Unset, "Unknown trace status")
	}

	return nil
}

// processSpans processes all spans in the trace
func (je *jaegerExporterComplete) processSpans(ctx context.Context, parentSpan oteltrace.Span, trace *domain.Trace) error {
	for i, span := range trace.Spans {
		// Create child span for each span processing
		spanName := fmt.Sprintf("process-span-%d", i)
		ctx, childSpan := je.tracer.Start(ctx, spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
			oteltrace.WithAttributes(
				attribute.String("span.id", string(span.ID)),
				attribute.String("span.service", string(span.Service)),
				attribute.String("span.operation", string(span.Operation)),
				attribute.String("span.status", string(span.Status)),
				attribute.Int64("span.duration_ns", span.Duration.Nanoseconds()),
			),
		)

		// Add span tags as attributes
		for key, value := range span.Tags {
			childSpan.SetAttributes(attribute.String("span.tag."+key, value))
		}

		// Process span logs
		if err := je.processSpanLogs(ctx, childSpan, span); err != nil {
			childSpan.SetStatus(codes.Error, err.Error())
			childSpan.End()
			return fmt.Errorf("failed to process span logs: %w", err)
		}

		// Set span status
		switch span.Status {
		case domain.SpanStatusOK:
			childSpan.SetStatus(codes.Ok, "Span processed successfully")
		case domain.SpanStatusError:
			childSpan.SetStatus(codes.Error, "Span contains errors")
		default:
			childSpan.SetStatus(codes.Unset, "Unknown span status")
		}

		childSpan.End()
	}

	return nil
}

// processSpanLogs processes logs within a span
func (je *jaegerExporterComplete) processSpanLogs(ctx context.Context, span oteltrace.Span, domainSpan domain.Span) error {
	for i, log := range domainSpan.Logs {
		// Add log as span event
		span.AddEvent("log",
			oteltrace.WithAttributes(
				attribute.String("log.message", log.Message),
				attribute.String("log.timestamp", log.Timestamp.Format(time.RFC3339)),
				attribute.Int("log.index", i),
			),
		)

		// Add log fields as attributes
		for key, value := range log.Fields {
			span.SetAttributes(attribute.String("log.field."+key, value))
		}
	}

	return nil
}

// recordTraceMetrics records metrics for the trace
func (je *jaegerExporterComplete) recordTraceMetrics(ctx context.Context, trace *domain.Trace) {
	// Record trace metrics
	je.telemetryManager.RecordMetric(ctx, "traces_exported_total", 1,
		attribute.String("service", string(trace.Service)),
		attribute.String("operation", string(trace.Operation)),
		attribute.String("status", string(trace.Status)),
	)

	// Record trace duration
	je.telemetryManager.RecordHistogram(ctx, "trace_duration_seconds", trace.Duration.Seconds(),
		attribute.String("service", string(trace.Service)),
		attribute.String("operation", string(trace.Operation)),
	)

	// Record spans count
	je.telemetryManager.RecordMetric(ctx, "spans_processed_total", float64(len(trace.Spans)),
		attribute.String("service", string(trace.Service)),
		attribute.String("operation", string(trace.Operation)),
	)
}

// validateTrace validates a trace before export
func (je *jaegerExporterComplete) validateTrace(trace *domain.Trace) error {
	if trace == nil {
		return fmt.Errorf("trace cannot be nil")
	}
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
