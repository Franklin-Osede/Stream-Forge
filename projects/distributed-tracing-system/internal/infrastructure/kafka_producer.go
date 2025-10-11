package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/segmentio/kafka-go"
)

// kafkaProducer implements the KafkaProducer interface
type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (domain.KafkaProducer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers list cannot be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	}

	return &kafkaProducer{
		writer: writer,
		topic:  topic,
	}, nil
}

// PublishTraceEvent publishes a trace event to Kafka
func (kp *kafkaProducer) PublishTraceEvent(ctx context.Context, trace *domain.Trace) error {
	// Validate trace
	if err := kp.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
	}

	// Convert trace to JSON
	traceData, err := json.Marshal(trace)
	if err != nil {
		return fmt.Errorf("failed to marshal trace: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(string(trace.ID)),
		Value: traceData,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{
				Key:   "service",
				Value: []byte(string(trace.Service)),
			},
			{
				Key:   "operation",
				Value: []byte(string(trace.Operation)),
			},
			{
				Key:   "status",
				Value: []byte(string(trace.Status)),
			},
		},
	}

	// Publish message
	if err := kp.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to publish trace event: %w", err)
	}

	return nil
}

// validateTrace validates a trace before publishing
func (kp *kafkaProducer) validateTrace(trace *domain.Trace) error {
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
	return nil
}

// Close closes the Kafka producer
func (kp *kafkaProducer) Close() error {
	return kp.writer.Close()
}
