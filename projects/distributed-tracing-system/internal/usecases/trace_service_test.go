package usecases

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockTraceRepository struct {
	mock.Mock
}

func (m *MockTraceRepository) Save(ctx context.Context, trace *domain.Trace) error {
	args := m.Called(ctx, trace)
	return args.Error(0)
}

func (m *MockTraceRepository) FindByID(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Trace), args.Error(1)
}

func (m *MockTraceRepository) Search(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Trace, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).([]*domain.Trace), args.Error(1)
}

func (m *MockTraceRepository) GetServices(ctx context.Context) ([]domain.ServiceName, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ServiceName), args.Error(1)
}

func (m *MockTraceRepository) GetOperations(ctx context.Context, service domain.ServiceName) ([]domain.OperationName, error) {
	args := m.Called(ctx, service)
	return args.Get(0).([]domain.OperationName), args.Error(1)
}

type MockPrometheusExporter struct {
	mock.Mock
}

func (m *MockPrometheusExporter) RecordTraceMetrics(trace *domain.Trace) error {
	args := m.Called(trace)
	return args.Error(0)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) PublishTraceEvent(ctx context.Context, trace *domain.Trace) error {
	args := m.Called(ctx, trace)
	return args.Error(0)
}

func TestTraceService_ProcessTrace_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	trace := &domain.Trace{
		ID:        "1234567890abcdef",
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Spans:     []domain.Span{},
		Tags:      map[string]string{},
		Status:    domain.TraceStatusSuccess,
	}

	// Setup expectations
	mockRepo.On("Save", ctx, trace).Return(nil)
	mockPrometheus.On("RecordTraceMetrics", trace).Return(nil)
	mockKafka.On("PublishTraceEvent", ctx, trace).Return(nil)

	// Act
	err := service.ProcessTrace(ctx, trace)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPrometheus.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestTraceService_ProcessTrace_InvalidTrace(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	invalidTrace := &domain.Trace{
		ID:        "", // Invalid: empty ID
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
	}

	// Act
	err := service.ProcessTrace(ctx, invalidTrace)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid trace")
	mockRepo.AssertNotCalled(t, "Save")
	mockPrometheus.AssertNotCalled(t, "RecordTraceMetrics")
	mockKafka.AssertNotCalled(t, "PublishTraceEvent")
}

func TestTraceService_ProcessTrace_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	trace := &domain.Trace{
		ID:        "1234567890abcdef",
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Spans:     []domain.Span{},
		Tags:      map[string]string{},
		Status:    domain.TraceStatusSuccess,
	}

	// Setup expectations
	expectedError := errors.New("database connection failed")
	mockRepo.On("Save", ctx, trace).Return(expectedError)

	// Act
	err := service.ProcessTrace(ctx, trace)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save trace")
	assert.Contains(t, err.Error(), "database connection failed")
	mockRepo.AssertExpectations(t)
	mockPrometheus.AssertNotCalled(t, "RecordTraceMetrics")
	mockKafka.AssertNotCalled(t, "PublishTraceEvent")
}

func TestTraceService_SearchTraces_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	criteria := &domain.SearchCriteria{
		Service: stringPtr("test-service"),
		Limit:   10,
		Offset:  0,
	}

	expectedTraces := []*domain.Trace{
		{
			ID:        "trace1",
			Service:   "test-service",
			Operation: "operation1",
			StartTime: time.Now().Add(-time.Hour),
			EndTime:   time.Now().Add(-time.Hour + time.Minute),
			Status:    domain.TraceStatusSuccess,
		},
		{
			ID:        "trace2",
			Service:   "test-service",
			Operation: "operation2",
			StartTime: time.Now().Add(-time.Hour),
			EndTime:   time.Now().Add(-time.Hour + time.Minute),
			Status:    domain.TraceStatusError,
		},
	}

	// Setup expectations
	mockRepo.On("Search", ctx, criteria).Return(expectedTraces, nil)

	// Act
	traces, err := service.SearchTraces(ctx, criteria)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, traces, 2)
	assert.Equal(t, expectedTraces, traces)
	mockRepo.AssertExpectations(t)
}

func TestTraceService_GetTrace_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	traceID := domain.TraceID("1234567890abcdef")
	expectedTrace := &domain.Trace{
		ID:        traceID,
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: time.Now().Add(-time.Second),
		EndTime:   time.Now(),
		Status:    domain.TraceStatusSuccess,
	}

	// Setup expectations
	mockRepo.On("FindByID", ctx, traceID).Return(expectedTrace, nil)

	// Act
	trace, err := service.GetTrace(ctx, traceID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTrace, trace)
	mockRepo.AssertExpectations(t)
}

func TestTraceService_GetServices_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	expectedServices := []domain.ServiceName{
		"service1",
		"service2",
		"service3",
	}

	// Setup expectations
	mockRepo.On("GetServices", ctx).Return(expectedServices, nil)

	// Act
	services, err := service.GetServices(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedServices, services)
	mockRepo.AssertExpectations(t)
}

func TestTraceService_GetOperations_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTraceRepository)
	mockPrometheus := new(MockPrometheusExporter)
	mockKafka := new(MockKafkaProducer)

	service := NewTraceService(mockRepo, mockPrometheus, mockKafka)

	ctx := context.Background()
	serviceName := domain.ServiceName("test-service")
	expectedOperations := []domain.OperationName{
		"operation1",
		"operation2",
		"operation3",
	}

	// Setup expectations
	mockRepo.On("GetOperations", ctx, serviceName).Return(expectedOperations, nil)

	// Act
	operations, err := service.GetOperations(ctx, serviceName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOperations, operations)
	mockRepo.AssertExpectations(t)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
