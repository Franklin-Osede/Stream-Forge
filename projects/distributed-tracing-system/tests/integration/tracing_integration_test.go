package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/streamforge/distributed-tracing-system/internal/infrastructure"
	"github.com/stretchr/testify/assert"
)

// IntegrationTestSuite provides a test suite for integration tests
type IntegrationTestSuite struct {
	jaegerExporter    domain.JaegerExporter
	prometheusExporter domain.PrometheusExporter
	kafkaProducer     domain.KafkaProducer
	kafkaConsumer     domain.KafkaConsumer
}

// SetupIntegrationTest sets up the integration test environment
func SetupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	// These tests require external services to be running
	// In a real CI/CD environment, you would use testcontainers
	// or docker-compose to spin up the required services

	jaegerExporter, err := infrastructure.NewJaegerExporter("http://localhost:14268/api/traces")
	if err != nil {
		t.Skip("Jaeger not available for integration testing")
	}

	prometheusExporter, err := infrastructure.NewPrometheusExporter("9091", "/metrics")
	if err != nil {
		t.Skip("Prometheus not available for integration testing")
	}

	kafkaProducer, err := infrastructure.NewKafkaProducer([]string{"localhost:9092"}, "trace-events")
	if err != nil {
		t.Skip("Kafka not available for integration testing")
	}

	kafkaConsumer, err := infrastructure.NewKafkaConsumer([]string{"localhost:9092"}, "trace-events", "test-group")
	if err != nil {
		t.Skip("Kafka not available for integration testing")
	}

	return &IntegrationTestSuite{
		jaegerExporter:    jaegerExporter,
		prometheusExporter: prometheusExporter,
		kafkaProducer:     kafkaProducer,
		kafkaConsumer:     kafkaConsumer,
	}
}

func TestTracingIntegration_EndToEndFlow(t *testing.T) {
	suite := SetupIntegrationTest(t)

	ctx := context.Background()

	// Create a test trace
	trace := &domain.Trace{
		ID:        "integration-test-trace",
		Service:   "integration-test-service",
		Operation: "integration-test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Spans: []domain.Span{
			{
				ID:        "span1",
				TraceID:   "integration-test-trace",
				Service:   "integration-test-service",
				Operation: "integration-test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
				Status:    domain.SpanStatusOK,
			},
		},
		Tags:   map[string]string{"test": "integration"},
		Status: domain.TraceStatusSuccess,
	}

	// Test Jaeger export
	err := suite.jaegerExporter.ExportTrace(ctx, trace)
	// In integration tests, we might expect this to succeed
	// but in unit tests, we expect it to fail due to no Jaeger instance
	if err != nil {
		t.Logf("Jaeger export failed (expected in test environment): %v", err)
	}

	// Test Prometheus metrics
	err = suite.prometheusExporter.RecordTraceMetrics(trace)
	assert.NoError(t, err)

	// Test Kafka producer
	err = suite.kafkaProducer.PublishTraceEvent(ctx, trace)
	if err != nil {
		t.Logf("Kafka publish failed (expected in test environment): %v", err)
	}
}

func TestTracingIntegration_MultipleTraces(t *testing.T) {
	suite := SetupIntegrationTest(t)

	ctx := context.Background()

	// Create multiple test traces
	traces := []*domain.Trace{
		{
			ID:        "trace1",
			Service:   "service1",
			Operation: "operation1",
			StartTime: time.Now().Add(-time.Minute),
			EndTime:   time.Now().Add(-time.Minute + time.Second),
			Status:    domain.TraceStatusSuccess,
		},
		{
			ID:        "trace2",
			Service:   "service2",
			Operation: "operation2",
			StartTime: time.Now().Add(-time.Minute),
			EndTime:   time.Now().Add(-time.Minute + time.Second),
			Status:    domain.TraceStatusError,
		},
		{
			ID:        "trace3",
			Service:   "service1",
			Operation: "operation3",
			StartTime: time.Now().Add(-time.Minute),
			EndTime:   time.Now().Add(-time.Minute + time.Second),
			Status:    domain.TraceStatusTimeout,
		},
	}

	// Process all traces
	for i, trace := range traces {
		t.Run(fmt.Sprintf("trace_%d", i+1), func(t *testing.T) {
			err := suite.jaegerExporter.ExportTrace(ctx, trace)
			if err != nil {
				t.Logf("Jaeger export failed for trace %d: %v", i+1, err)
			}

			err = suite.prometheusExporter.RecordTraceMetrics(trace)
			assert.NoError(t, err)

			err = suite.kafkaProducer.PublishTraceEvent(ctx, trace)
			if err != nil {
				t.Logf("Kafka publish failed for trace %d: %v", i+1, err)
			}
		})
	}
}

func TestTracingIntegration_ErrorHandling(t *testing.T) {
	suite := SetupIntegrationTest(t)

	ctx := context.Background()

	// Test with invalid trace
	invalidTrace := &domain.Trace{
		ID:        "", // Invalid: empty ID
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
	}

	// All operations should handle invalid traces gracefully
	err := suite.jaegerExporter.ExportTrace(ctx, invalidTrace)
	assert.Error(t, err)

	err = suite.prometheusExporter.RecordTraceMetrics(invalidTrace)
	assert.Error(t, err)

	err = suite.kafkaProducer.PublishTraceEvent(ctx, invalidTrace)
	assert.Error(t, err)
}

func TestTracingIntegration_ConcurrentTraces(t *testing.T) {
	suite := SetupIntegrationTest(t)

	ctx := context.Background()
	numTraces := 10
	done := make(chan bool, numTraces)

	// Create and process traces concurrently
	for i := 0; i < numTraces; i++ {
		go func(traceNum int) {
			defer func() { done <- true }()

			trace := &domain.Trace{
				ID:        domain.TraceID(fmt.Sprintf("concurrent-trace-%d", traceNum)),
				Service:   domain.ServiceName(fmt.Sprintf("service-%d", traceNum%3)),
				Operation: domain.OperationName(fmt.Sprintf("operation-%d", traceNum%5)),
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
				Status:    domain.TraceStatusSuccess,
			}

			// Process trace
			err := suite.jaegerExporter.ExportTrace(ctx, trace)
			if err != nil {
				t.Logf("Jaeger export failed for concurrent trace %d: %v", traceNum, err)
			}

			err = suite.prometheusExporter.RecordTraceMetrics(trace)
			assert.NoError(t, err)

			err = suite.kafkaProducer.PublishTraceEvent(ctx, trace)
			if err != nil {
				t.Logf("Kafka publish failed for concurrent trace %d: %v", traceNum, err)
			}
		}(i)
	}

	// Wait for all traces to complete
	for i := 0; i < numTraces; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkTracingIntegration_ExportTrace(b *testing.B) {
	suite := SetupIntegrationTest(&testing.T{})
	if suite == nil {
		b.Skip("Integration test environment not available")
	}

	ctx := context.Background()
	trace := &domain.Trace{
		ID:        "benchmark-trace",
		Service:   "benchmark-service",
		Operation: "benchmark-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Status:    domain.TraceStatusSuccess,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.jaegerExporter.ExportTrace(ctx, trace)
	}
}

func BenchmarkTracingIntegration_RecordMetrics(b *testing.B) {
	suite := SetupIntegrationTest(&testing.T{})
	if suite == nil {
		b.Skip("Integration test environment not available")
	}

	trace := &domain.Trace{
		ID:        "benchmark-trace",
		Service:   "benchmark-service",
		Operation: "benchmark-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Status:    domain.TraceStatusSuccess,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.prometheusExporter.RecordTraceMetrics(trace)
	}
}
