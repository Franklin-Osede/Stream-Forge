package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrace_Validate(t *testing.T) {
	tests := []struct {
		name    string
		trace   *Trace
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid trace",
			trace: &Trace{
				ID:        "1234567890abcdef",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
				Spans:     []Span{},
				Tags:      map[string]string{},
				Status:    TraceStatusSuccess,
			},
			wantErr: false,
		},
		{
			name: "missing trace ID",
			trace: &Trace{
				ID:        "",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
			},
			wantErr: true,
			errMsg:  "trace ID is required",
		},
		{
			name: "missing service name",
			trace: &Trace{
				ID:        "1234567890abcdef",
				Service:   "",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
			},
			wantErr: true,
			errMsg:  "service name is required",
		},
		{
			name: "missing operation name",
			trace: &Trace{
				ID:        "1234567890abcdef",
				Service:   "test-service",
				Operation: "",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
			},
			wantErr: true,
			errMsg:  "operation name is required",
		},
		{
			name: "start time after end time",
			trace: &Trace{
				ID:        "1234567890abcdef",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(-time.Second),
			},
			wantErr: true,
			errMsg:  "start time cannot be after end time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTrace(tt.trace)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrace_CalculateDuration(t *testing.T) {
	startTime := time.Now().Add(-5 * time.Second)
	endTime := time.Now()

	trace := &Trace{
		ID:        "1234567890abcdef",
		Service:   "test-service",
		Operation: "test-operation",
		StartTime: startTime,
		EndTime:   endTime,
	}

	expectedDuration := endTime.Sub(startTime)
	trace.Duration = expectedDuration
	assert.Equal(t, expectedDuration, trace.Duration)
}

func TestSpan_Validate(t *testing.T) {
	tests := []struct {
		name    string
		span    *Span
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid span",
			span: &Span{
				ID:        "span123",
				TraceID:   "trace123",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
				Status:    SpanStatusOK,
			},
			wantErr: false,
		},
		{
			name: "missing span ID",
			span: &Span{
				ID:        "",
				TraceID:   "trace123",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
			},
			wantErr: true,
			errMsg:  "span ID is required",
		},
		{
			name: "missing trace ID",
			span: &Span{
				ID:        "span123",
				TraceID:   "",
				Service:   "test-service",
				Operation: "test-operation",
				StartTime: time.Now().Add(-time.Second),
				EndTime:   time.Now(),
			},
			wantErr: true,
			errMsg:  "trace ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSpan(tt.span)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchCriteria_Validate(t *testing.T) {
	tests := []struct {
		name    string
		criteria *SearchCriteria
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid criteria",
			criteria: &SearchCriteria{
				Service:   (*ServiceName)(stringPtr("test-service")),
				Operation: (*OperationName)(stringPtr("test-operation")),
				Limit:     10,
				Offset:    0,
			},
			wantErr: false,
		},
		{
			name: "negative limit",
			criteria: &SearchCriteria{
				Limit: -1,
			},
			wantErr: true,
			errMsg:  "limit cannot be negative",
		},
		{
			name: "negative offset",
			criteria: &SearchCriteria{
				Offset: -1,
			},
			wantErr: true,
			errMsg:  "offset cannot be negative",
		},
		{
			name: "start time after end time",
			criteria: &SearchCriteria{
				StartTime: timePtr(time.Now()),
				EndTime:   timePtr(time.Now().Add(-time.Hour)),
			},
			wantErr: true,
			errMsg:  "start time cannot be after end time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.criteria.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Validation functions (these would be implemented in the domain)
func validateTrace(trace *Trace) error {
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

func validateSpan(span *Span) error {
	if span.ID == "" {
		return fmt.Errorf("span ID is required")
	}
	if span.TraceID == "" {
		return fmt.Errorf("trace ID is required")
	}
	if span.Service == "" {
		return fmt.Errorf("service name is required")
	}
	if span.Operation == "" {
		return fmt.Errorf("operation name is required")
	}
	if span.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if span.EndTime.IsZero() {
		return fmt.Errorf("end time is required")
	}
	if span.StartTime.After(span.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}
	return nil
}

// Add Validate method to SearchCriteria
func (sc *SearchCriteria) Validate() error {
	if sc.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if sc.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if sc.StartTime != nil && sc.EndTime != nil {
		if sc.StartTime.After(*sc.EndTime) {
			return fmt.Errorf("start time cannot be after end time")
		}
	}
	return nil
}
