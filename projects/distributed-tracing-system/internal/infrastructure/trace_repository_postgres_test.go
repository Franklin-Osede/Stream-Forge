package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockJaegerExporter for testing
type MockJaegerExporter struct{}

func (m *MockJaegerExporter) ExportTrace(ctx context.Context, trace *domain.Trace) error {
	return nil
}

func TestTraceRepositoryPostgres_NewTraceRepositoryPostgres(t *testing.T) {
	// This test requires a real PostgreSQL database
	// In a real test environment, you would use testcontainers or a test database
	t.Skip("Skipping test that requires PostgreSQL database")
}

func TestTraceRepositoryPostgres_Save(t *testing.T) {
	// This test requires a real PostgreSQL database
	t.Skip("Skipping test that requires PostgreSQL database")
}

func TestTraceRepositoryPostgres_FindByID(t *testing.T) {
	// This test requires a real PostgreSQL database
	t.Skip("Skipping test that requires PostgreSQL database")
}

func TestTraceRepositoryPostgres_Search(t *testing.T) {
	// This test requires a real PostgreSQL database
	t.Skip("Skipping test that requires PostgreSQL database")
}

func TestTraceRepositoryPostgres_GetServices(t *testing.T) {
	// This test requires a real PostgreSQL database
	t.Skip("Skipping test that requires PostgreSQL database")
}

func TestTraceRepositoryPostgres_GetOperations(t *testing.T) {
	// This test requires a real PostgreSQL database
	t.Skip("Skipping test that requires PostgreSQL database")
}

// Integration test that can be run with a real database
func TestTraceRepositoryPostgres_Integration(t *testing.T) {
	// Skip if no database connection string is provided
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=tracing_system sslmode=disable"
	
	// Create mock Jaeger exporter
	mockJaeger := &MockJaegerExporter{}
	
	// Create repository
	repo, err := NewTraceRepositoryPostgres(dsn, mockJaeger)
	if err != nil {
		t.Skip("Skipping integration test - database not available")
	}
	defer func() {
		if pgRepo, ok := repo.(*traceRepositoryPostgres); ok {
			pgRepo.Close()
		}
	}()

	ctx := context.Background()

	// Test data
	trace := &domain.Trace{
		ID:        "test-trace-123",
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Minute),
		EndTime:   time.Now(),
		Duration:  time.Minute,
		Spans: []domain.Span{
			{
				ID:        "span-1",
				TraceID:   "test-trace-123",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Minute),
				EndTime:   time.Now(),
				Duration:  time.Minute,
				Status:    domain.SpanStatusOK,
				Tags:      map[string]string{"test": "true"},
				Logs: []domain.Log{
					{
						Timestamp: time.Now(),
						Message:   "Test log message",
						Fields:    map[string]string{"level": "info"},
					},
				},
			},
		},
		Tags:   map[string]string{"environment": "test"},
		Status: domain.TraceStatusSuccess,
	}

	// Test Save
	err = repo.Save(ctx, trace)
	require.NoError(t, err)

	// Test FindByID
	foundTrace, err := repo.FindByID(ctx, trace.ID)
	require.NoError(t, err)
	assert.Equal(t, trace.ID, foundTrace.ID)
	assert.Equal(t, trace.Service, foundTrace.Service)
	assert.Equal(t, trace.Operation, foundTrace.Operation)
	assert.Equal(t, trace.Status, foundTrace.Status)
	assert.Len(t, foundTrace.Spans, 1)

	// Test Search
	criteria := &domain.SearchCriteria{
		Service: &trace.Service,
		Limit:   10,
		Offset:  0,
	}
	
	searchResults, err := repo.Search(ctx, criteria)
	require.NoError(t, err)
	assert.Len(t, searchResults, 1)
	assert.Equal(t, trace.ID, searchResults[0].ID)

	// Test GetServices
	services, err := repo.GetServices(ctx)
	require.NoError(t, err)
	assert.Contains(t, services, trace.Service)

	// Test GetOperations
	operations, err := repo.GetOperations(ctx, trace.Service)
	require.NoError(t, err)
	assert.Contains(t, operations, trace.Operation)
}
