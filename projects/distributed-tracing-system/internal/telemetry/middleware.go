package telemetry

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware provides OpenTelemetry instrumentation for HTTP requests
type HTTPMiddleware struct {
	tracer trace.Tracer
	meter  metric.Meter
}

// NewHTTPMiddleware creates a new HTTP middleware
func NewHTTPMiddleware(tracer trace.Tracer, meter metric.Meter) *HTTPMiddleware {
	return &HTTPMiddleware{
		tracer: tracer,
		meter:  meter,
	}
}

// GinMiddleware returns a Gin middleware for OpenTelemetry instrumentation
func (m *HTTPMiddleware) GinMiddleware() gin.HandlerFunc {
	// Create metrics
	requestCounter, _ := m.meter.Int64Counter("http_requests_total")
	requestDuration, _ := m.meter.Float64Histogram("http_request_duration_seconds")
	requestSize, _ := m.meter.Int64Histogram("http_request_size_bytes")
	responseSize, _ := m.meter.Int64Histogram("http_response_size_bytes")

	return func(c *gin.Context) {
		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		
		// Start span
		spanName := c.Request.Method + " " + c.FullPath()
		ctx, span := m.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("http.route", c.FullPath()),
			),
		)

		// Set context
		c.Request = c.Request.WithContext(ctx)

		// Record request size
		requestSize.Record(ctx, int64(c.Request.ContentLength), metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
		))

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get response size
		responseSizeBytes := c.Writer.Size()

		// Set span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int64("http.response.size", int64(responseSizeBytes)),
			attribute.Float64("http.request.duration", duration),
		)

		// Set span status
		if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, http.StatusText(c.Writer.Status()))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Record metrics
		requestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.status_code", strconv.Itoa(c.Writer.Status())),
		))

		requestDuration.Record(ctx, duration, metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.status_code", strconv.Itoa(c.Writer.Status())),
		))

		responseSize.Record(ctx, int64(responseSizeBytes), metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
		))

		// End span
		span.End()
	}
}

// DatabaseMiddleware provides OpenTelemetry instrumentation for database operations
type DatabaseMiddleware struct {
	tracer trace.Tracer
	meter  metric.Meter
}

// NewDatabaseMiddleware creates a new database middleware
func NewDatabaseMiddleware(tracer trace.Tracer, meter metric.Meter) *DatabaseMiddleware {
	return &DatabaseMiddleware{
		tracer: tracer,
		meter:  meter,
	}
}

// InstrumentQuery instruments a database query
func (m *DatabaseMiddleware) InstrumentQuery(ctx context.Context, operation string, query string) (context.Context, trace.Span) {
	ctx, span := m.tracer.Start(ctx, "db."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.operation", operation),
			attribute.String("db.statement", query),
			attribute.String("db.system", "postgresql"),
		),
	)

	// Record database metrics
	queryCounter, _ := m.meter.Int64Counter("db_queries_total")
	queryCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("db.operation", operation),
		attribute.String("db.system", "postgresql"),
	))

	return ctx, span
}

// InstrumentTransaction instruments a database transaction
func (m *DatabaseMiddleware) InstrumentTransaction(ctx context.Context, operation string) (context.Context, trace.Span) {
	ctx, span := m.tracer.Start(ctx, "db.transaction."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.operation", "transaction"),
			attribute.String("db.transaction.name", operation),
			attribute.String("db.system", "postgresql"),
		),
	)

	return ctx, span
}

// KafkaMiddleware provides OpenTelemetry instrumentation for Kafka operations
type KafkaMiddleware struct {
	tracer trace.Tracer
	meter  metric.Meter
}

// NewKafkaMiddleware creates a new Kafka middleware
func NewKafkaMiddleware(tracer trace.Tracer, meter metric.Meter) *KafkaMiddleware {
	return &KafkaMiddleware{
		tracer: tracer,
		meter:  meter,
	}
}

// InstrumentProducer instruments a Kafka producer operation
func (m *KafkaMiddleware) InstrumentProducer(ctx context.Context, topic string, key string) (context.Context, trace.Span) {
	ctx, span := m.tracer.Start(ctx, "kafka.produce",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.destination_kind", "topic"),
			attribute.String("messaging.message_key", key),
		),
	)

	// Record Kafka metrics
	producerCounter, _ := m.meter.Int64Counter("kafka_messages_produced_total")
	producerCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("kafka.topic", topic),
		attribute.String("kafka.key", key),
	))

	return ctx, span
}

// InstrumentConsumer instruments a Kafka consumer operation
func (m *KafkaMiddleware) InstrumentConsumer(ctx context.Context, topic string, groupID string) (context.Context, trace.Span) {
	ctx, span := m.tracer.Start(ctx, "kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.destination_kind", "topic"),
			attribute.String("messaging.consumer_group", groupID),
		),
	)

	// Record Kafka metrics
	consumerCounter, _ := m.meter.Int64Counter("kafka_messages_consumed_total")
	consumerCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("kafka.topic", topic),
		attribute.String("kafka.group_id", groupID),
	))

	return ctx, span
}
