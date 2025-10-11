package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/segmentio/kafka-go"
)

// kafkaConsumer implements the KafkaConsumer interface
type kafkaConsumer struct {
	reader *kafka.Reader
	topic  string
	groupID string
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic, groupID string) (domain.KafkaConsumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers list cannot be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}
	if groupID == "" {
		return nil, fmt.Errorf("group ID cannot be empty")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &kafkaConsumer{
		reader:  reader,
		topic:   topic,
		groupID: groupID,
	}, nil
}

// Start starts the Kafka consumer
func (kc *kafkaConsumer) Start(ctx context.Context, traceService domain.TraceService) error {
	log.Printf("Starting Kafka consumer for topic: %s, group: %s", kc.topic, kc.groupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping due to context cancellation")
			return ctx.Err()
		default:
			// Read message with timeout
			message, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}

			// Process message
			if err := kc.processMessage(ctx, message, traceService); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue processing other messages
				continue
			}
		}
	}
}

// processMessage processes a single Kafka message
func (kc *kafkaConsumer) processMessage(ctx context.Context, message kafka.Message, traceService domain.TraceService) error {
	// Parse trace from message
	var trace domain.Trace
	if err := json.Unmarshal(message.Value, &trace); err != nil {
		return fmt.Errorf("failed to unmarshal trace: %w", err)
	}

	// Validate trace
	if err := kc.validateTrace(&trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Process trace through service
	if err := traceService.ProcessTrace(ctx, &trace); err != nil {
		return fmt.Errorf("failed to process trace: %w", err)
	}

	log.Printf("Successfully processed trace: %s", trace.ID)
	return nil
}

// validateTrace validates a trace from Kafka message
func (kc *kafkaConsumer) validateTrace(trace *domain.Trace) error {
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

// Close closes the Kafka consumer
func (kc *kafkaConsumer) Close() error {
	return kc.reader.Close()
}

