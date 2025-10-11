package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
)

// traceRepository implements the TraceRepository interface
type traceRepository struct {
	jaegerExporter domain.JaegerExporter
	// In a real implementation, you would have a database connection here
	// db *sql.DB
}

// NewTraceRepository creates a new trace repository
func NewTraceRepository(jaegerExporter domain.JaegerExporter) domain.TraceRepository {
	return &traceRepository{
		jaegerExporter: jaegerExporter,
	}
}

// Save saves a trace to the repository
func (tr *traceRepository) Save(ctx context.Context, trace *domain.Trace) error {
	// Validate trace
	if err := tr.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Export to Jaeger
	if err := tr.jaegerExporter.ExportTrace(ctx, trace); err != nil {
		return fmt.Errorf("failed to export trace to Jaeger: %w", err)
	}

	// In a real implementation, you would also save to database
	// if err := tr.saveToDatabase(ctx, trace); err != nil {
	//     return fmt.Errorf("failed to save trace to database: %w", err)
	// }

	return nil
}

// FindByID finds a trace by ID
func (tr *traceRepository) FindByID(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	if id == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	// In a real implementation, you would query the database
	// For now, return a mock trace
	return &domain.Trace{
		ID:        id,
		Service:   "mock-service",
		Operation: "mock-operation",
		StartTime: time.Now().Add(-time.Minute),
		EndTime:   time.Now(),
		Duration:  time.Minute,
		Status:    domain.TraceStatusSuccess,
	}, nil
}

// Search searches for traces based on criteria
func (tr *traceRepository) Search(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Trace, error) {
	// Validate criteria
	if err := tr.validateSearchCriteria(criteria); err != nil {
		return nil, fmt.Errorf("invalid search criteria: %w", err)
	}

	// In a real implementation, you would query the database
	// For now, return mock traces
	mockTraces := []*domain.Trace{
		{
			ID:        "trace1",
			Service:   "service1",
			Operation: "operation1",
			StartTime: time.Now().Add(-time.Hour),
			EndTime:   time.Now().Add(-time.Hour + time.Minute),
			Duration:  time.Minute,
			Status:    domain.TraceStatusSuccess,
		},
		{
			ID:        "trace2",
			Service:   "service2",
			Operation: "operation2",
			StartTime: time.Now().Add(-time.Hour),
			EndTime:   time.Now().Add(-time.Hour + time.Minute),
			Duration:  time.Minute,
			Status:    domain.TraceStatusError,
		},
	}

	// Apply limit and offset
	start := criteria.Offset
	end := start + criteria.Limit
	if end > len(mockTraces) {
		end = len(mockTraces)
	}
	if start >= len(mockTraces) {
		return []*domain.Trace{}, nil
	}

	return mockTraces[start:end], nil
}

// GetServices returns all available services
func (tr *traceRepository) GetServices(ctx context.Context) ([]domain.ServiceName, error) {
	// In a real implementation, you would query the database
	// For now, return mock services
	return []domain.ServiceName{
		"service1",
		"service2",
		"service3",
	}, nil
}

// GetOperations returns all operations for a specific service
func (tr *traceRepository) GetOperations(ctx context.Context, service domain.ServiceName) ([]domain.OperationName, error) {
	if service == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// In a real implementation, you would query the database
	// For now, return mock operations
	return []domain.OperationName{
		"operation1",
		"operation2",
		"operation3",
	}, nil
}

// validateTrace validates a trace before saving
func (tr *traceRepository) validateTrace(trace *domain.Trace) error {
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

// validateSearchCriteria validates search criteria
func (tr *traceRepository) validateSearchCriteria(criteria *domain.SearchCriteria) error {
	if criteria == nil {
		return fmt.Errorf("search criteria cannot be nil")
	}
	if criteria.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if criteria.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if criteria.StartTime != nil && criteria.EndTime != nil {
		if criteria.StartTime.After(*criteria.EndTime) {
			return fmt.Errorf("start time cannot be after end time")
		}
	}
	return nil
}

