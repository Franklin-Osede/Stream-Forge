package domain

import (
	"context"
	"time"
)

// TraceID represents a unique trace identifier
type TraceID string

// SpanID represents a unique span identifier
type SpanID string

// ServiceName represents the name of a service
type ServiceName string

// OperationName represents the name of an operation
type OperationName string

// Trace represents a distributed trace
type Trace struct {
	ID         TraceID    `json:"id"`
	Service    ServiceName `json:"service"`
	Operation  OperationName `json:"operation"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	Spans      []Span     `json:"spans"`
	Tags       map[string]string `json:"tags"`
	Status     TraceStatus `json:"status"`
}

// Span represents a single operation within a trace
type Span struct {
	ID          SpanID     `json:"id"`
	TraceID     TraceID    `json:"trace_id"`
	ParentID    *SpanID    `json:"parent_id,omitempty"`
	Service     ServiceName `json:"service"`
	Operation   OperationName `json:"operation"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	Tags        map[string]string `json:"tags"`
	Logs        []Log      `json:"logs,omitempty"`
	Status      SpanStatus `json:"status"`
}

// Log represents a log entry within a span
type Log struct {
	Timestamp time.Time         `json:"timestamp"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
}

// TraceStatus represents the status of a trace
type TraceStatus string

const (
	TraceStatusSuccess TraceStatus = "success"
	TraceStatusError   TraceStatus = "error"
	TraceStatusTimeout TraceStatus = "timeout"
)

// SpanStatus represents the status of a span
type SpanStatus string

const (
	SpanStatusOK    SpanStatus = "ok"
	SpanStatusError SpanStatus = "error"
)

// TraceRepository defines the interface for trace persistence
type TraceRepository interface {
	Save(ctx context.Context, trace *Trace) error
	FindByID(ctx context.Context, id TraceID) (*Trace, error)
	Search(ctx context.Context, criteria *SearchCriteria) ([]*Trace, error)
	GetServices(ctx context.Context) ([]ServiceName, error)
	GetOperations(ctx context.Context, service ServiceName) ([]OperationName, error)
}

// SearchCriteria defines search parameters for traces
type SearchCriteria struct {
	Service     *ServiceName    `json:"service,omitempty"`
	Operation   *OperationName  `json:"operation,omitempty"`
	StartTime   *time.Time     `json:"start_time,omitempty"`
	EndTime     *time.Time     `json:"end_time,omitempty"`
	MinDuration *time.Duration `json:"min_duration,omitempty"`
	MaxDuration *time.Duration `json:"max_duration,omitempty"`
	Status      *TraceStatus   `json:"status,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Limit       int            `json:"limit,omitempty"`
	Offset      int            `json:"offset,omitempty"`
}

// TraceService defines the business logic for trace operations
type TraceService interface {
	ProcessTrace(ctx context.Context, trace *Trace) error
	SearchTraces(ctx context.Context, criteria *SearchCriteria) ([]*Trace, error)
	GetTrace(ctx context.Context, id TraceID) (*Trace, error)
	GetServices(ctx context.Context) ([]ServiceName, error)
	GetOperations(ctx context.Context, service ServiceName) ([]OperationName, error)
	GetMetrics(ctx context.Context) (*TraceMetrics, error)
}

// TraceMetrics represents aggregated metrics for traces
type TraceMetrics struct {
	TotalTraces     int64         `json:"total_traces"`
	TotalSpans      int64         `json:"total_spans"`
	AverageDuration time.Duration `json:"average_duration"`
	ErrorRate       float64       `json:"error_rate"`
	Throughput      float64       `json:"throughput"`
	Services        []ServiceMetrics `json:"services"`
}

// ServiceMetrics represents metrics for a specific service
type ServiceMetrics struct {
	Service         ServiceName   `json:"service"`
	TotalTraces     int64         `json:"total_traces"`
	AverageDuration time.Duration `json:"average_duration"`
	ErrorRate       float64       `json:"error_rate"`
	Throughput      float64       `json:"throughput"`
}
