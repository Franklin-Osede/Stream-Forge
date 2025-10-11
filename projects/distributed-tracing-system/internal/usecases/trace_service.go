package usecases

import (
	"context"
	"fmt"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
)

// traceService implements the TraceService interface
type traceService struct {
	repo            domain.TraceRepository
	prometheusExporter domain.PrometheusExporter
	kafkaProducer   domain.KafkaProducer
}

// NewTraceService creates a new trace service
func NewTraceService(
	repo domain.TraceRepository,
	prometheusExporter domain.PrometheusExporter,
	kafkaProducer domain.KafkaProducer,
) domain.TraceService {
	return &traceService{
		repo:            repo,
		prometheusExporter: prometheusExporter,
		kafkaProducer:   kafkaProducer,
	}
}

// ProcessTrace processes a new trace
func (s *traceService) ProcessTrace(ctx context.Context, trace *domain.Trace) error {
	// Validate trace
	if err := s.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Calculate duration if not set
	if trace.Duration == 0 {
		trace.Duration = trace.EndTime.Sub(trace.StartTime)
	}

	// Save trace to repository
	if err := s.repo.Save(ctx, trace); err != nil {
		return fmt.Errorf("failed to save trace: %w", err)
	}

	// Export metrics to Prometheus
	if err := s.prometheusExporter.RecordTraceMetrics(trace); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to record metrics: %v\n", err)
	}

	// Publish trace event to Kafka
	if err := s.kafkaProducer.PublishTraceEvent(ctx, trace); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish trace event: %v\n", err)
	}

	return nil
}

// SearchTraces searches for traces based on criteria
func (s *traceService) SearchTraces(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Trace, error) {
	return s.repo.Search(ctx, criteria)
}

// GetTrace retrieves a specific trace by ID
func (s *traceService) GetTrace(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	return s.repo.FindByID(ctx, id)
}

// GetServices retrieves all available services
func (s *traceService) GetServices(ctx context.Context) ([]domain.ServiceName, error) {
	return s.repo.GetServices(ctx)
}

// GetOperations retrieves all operations for a specific service
func (s *traceService) GetOperations(ctx context.Context, service domain.ServiceName) ([]domain.OperationName, error) {
	return s.repo.GetOperations(ctx, service)
}

// GetMetrics retrieves aggregated trace metrics
func (s *traceService) GetMetrics(ctx context.Context) (*domain.TraceMetrics, error) {
	// This would typically query Prometheus or calculate from the repository
	// For now, return a basic implementation
	return &domain.TraceMetrics{
		TotalTraces:     0,
		TotalSpans:      0,
		AverageDuration: 0,
		ErrorRate:       0.0,
		Throughput:      0.0,
		Services:       []domain.ServiceMetrics{},
	}, nil
}

// validateTrace validates a trace before processing
func (s *traceService) validateTrace(trace *domain.Trace) error {
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
