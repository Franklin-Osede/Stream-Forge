package domain

import "context"

// PrometheusExporter defines the interface for Prometheus metrics export
type PrometheusExporter interface {
	RecordTraceMetrics(trace *Trace) error
}

// KafkaProducer defines the interface for Kafka message publishing
type KafkaProducer interface {
	PublishTraceEvent(ctx context.Context, trace *Trace) error
}

// KafkaConsumer defines the interface for Kafka message consumption
type KafkaConsumer interface {
	Start(ctx context.Context, traceService TraceService) error
}

// JaegerExporter defines the interface for Jaeger trace export
type JaegerExporter interface {
	ExportTrace(ctx context.Context, trace *Trace) error
}
