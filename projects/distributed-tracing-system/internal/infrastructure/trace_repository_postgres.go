package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// traceRepositoryPostgres implements the TraceRepository interface with PostgreSQL
type traceRepositoryPostgres struct {
	db             *sqlx.DB
	jaegerExporter domain.JaegerExporter
}

// NewTraceRepositoryPostgres creates a new PostgreSQL trace repository
func NewTraceRepositoryPostgres(dsn string, jaegerExporter domain.JaegerExporter) (domain.TraceRepository, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &traceRepositoryPostgres{
		db:             db,
		jaegerExporter: jaegerExporter,
	}, nil
}

// createTables creates the necessary database tables
func createTables(db *sqlx.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS traces (
			id VARCHAR(255) PRIMARY KEY,
			service VARCHAR(255) NOT NULL,
			operation VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			duration BIGINT NOT NULL,
			status VARCHAR(50) NOT NULL,
			tags JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS spans (
			id VARCHAR(255) PRIMARY KEY,
			trace_id VARCHAR(255) NOT NULL,
			parent_id VARCHAR(255),
			service VARCHAR(255) NOT NULL,
			operation VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			duration BIGINT NOT NULL,
			status VARCHAR(50) NOT NULL,
			tags JSONB,
			logs JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (trace_id) REFERENCES traces(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_service ON traces(service)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_operation ON traces(operation)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_start_time ON traces(start_time)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_status ON traces(status)`,
		`CREATE INDEX IF NOT EXISTS idx_spans_trace_id ON spans(trace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_spans_service ON spans(service)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

// Save saves a trace to the repository
func (tr *traceRepositoryPostgres) Save(ctx context.Context, trace *domain.Trace) error {
	// Validate trace
	if err := tr.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Start transaction
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Save trace
	if err := tr.saveTrace(ctx, tx, trace); err != nil {
		return fmt.Errorf("failed to save trace: %w", err)
	}

	// Save spans
	if err := tr.saveSpans(ctx, tx, trace); err != nil {
		return fmt.Errorf("failed to save spans: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Export to Jaeger
	if err := tr.jaegerExporter.ExportTrace(ctx, trace); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to export trace to Jaeger: %v\n", err)
	}

	return nil
}

// saveTrace saves a trace to the database
func (tr *traceRepositoryPostgres) saveTrace(ctx context.Context, tx *sqlx.Tx, trace *domain.Trace) error {
	query := `
		INSERT INTO traces (id, service, operation, start_time, end_time, duration, status, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			service = EXCLUDED.service,
			operation = EXCLUDED.operation,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			duration = EXCLUDED.duration,
			status = EXCLUDED.status,
			tags = EXCLUDED.tags
	`

	tagsJSON, err := json.Marshal(trace.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = tx.ExecContext(ctx, query,
		trace.ID,
		trace.Service,
		trace.Operation,
		trace.StartTime,
		trace.EndTime,
		trace.Duration.Nanoseconds(),
		trace.Status,
		tagsJSON,
	)

	return err
}

// saveSpans saves spans to the database
func (tr *traceRepositoryPostgres) saveSpans(ctx context.Context, tx *sqlx.Tx, trace *domain.Trace) error {
	if len(trace.Spans) == 0 {
		return nil
	}

	query := `
		INSERT INTO spans (id, trace_id, parent_id, service, operation, start_time, end_time, duration, status, tags, logs)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			trace_id = EXCLUDED.trace_id,
			parent_id = EXCLUDED.parent_id,
			service = EXCLUDED.service,
			operation = EXCLUDED.operation,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			duration = EXCLUDED.duration,
			status = EXCLUDED.status,
			tags = EXCLUDED.tags,
			logs = EXCLUDED.logs
	`

	for _, span := range trace.Spans {
		tagsJSON, err := json.Marshal(span.Tags)
		if err != nil {
			return fmt.Errorf("failed to marshal span tags: %w", err)
		}

		logsJSON, err := json.Marshal(span.Logs)
		if err != nil {
			return fmt.Errorf("failed to marshal span logs: %w", err)
		}

		var parentID *string
		if span.ParentID != nil {
			parentIDStr := string(*span.ParentID)
			parentID = &parentIDStr
		}

		_, err = tx.ExecContext(ctx, query,
			span.ID,
			span.TraceID,
			parentID,
			span.Service,
			span.Operation,
			span.StartTime,
			span.EndTime,
			span.Duration.Nanoseconds(),
			span.Status,
			tagsJSON,
			logsJSON,
		)

		if err != nil {
			return fmt.Errorf("failed to save span %s: %w", span.ID, err)
		}
	}

	return nil
}

// FindByID finds a trace by ID
func (tr *traceRepositoryPostgres) FindByID(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	if id == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	// Query trace
	var trace domain.Trace
	var tagsJSON []byte
	query := `SELECT id, service, operation, start_time, end_time, duration, status, tags FROM traces WHERE id = $1`
	
	err := tr.db.QueryRowContext(ctx, query, id).Scan(
		&trace.ID,
		&trace.Service,
		&trace.Operation,
		&trace.StartTime,
		&trace.EndTime,
		&trace.Duration,
		&trace.Status,
		&tagsJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trace not found")
		}
		return nil, fmt.Errorf("failed to query trace: %w", err)
	}

	// Parse tags
	if err := json.Unmarshal(tagsJSON, &trace.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	// Query spans
	spans, err := tr.getSpansByTraceID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get spans: %w", err)
	}
	trace.Spans = spans

	return &trace, nil
}

// getSpansByTraceID gets all spans for a trace
func (tr *traceRepositoryPostgres) getSpansByTraceID(ctx context.Context, traceID domain.TraceID) ([]domain.Span, error) {
	query := `
		SELECT id, trace_id, parent_id, service, operation, start_time, end_time, duration, status, tags, logs
		FROM spans WHERE trace_id = $1 ORDER BY start_time
	`

	rows, err := tr.db.QueryContext(ctx, query, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query spans: %w", err)
	}
	defer rows.Close()

	var spans []domain.Span
	for rows.Next() {
		var span domain.Span
		var parentID *string
		var tagsJSON, logsJSON []byte

		err := rows.Scan(
			&span.ID,
			&span.TraceID,
			&parentID,
			&span.Service,
			&span.Operation,
			&span.StartTime,
			&span.EndTime,
			&span.Duration,
			&span.Status,
			&tagsJSON,
			&logsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan span: %w", err)
		}

		// Set parent ID
		if parentID != nil {
			parentIDStr := domain.SpanID(*parentID)
			span.ParentID = &parentIDStr
		}

		// Parse tags and logs
		if err := json.Unmarshal(tagsJSON, &span.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal span tags: %w", err)
		}

		if err := json.Unmarshal(logsJSON, &span.Logs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal span logs: %w", err)
		}

		spans = append(spans, span)
	}

	return spans, nil
}

// Search searches for traces based on criteria
func (tr *traceRepositoryPostgres) Search(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Trace, error) {
	// Validate criteria
	if err := tr.validateSearchCriteria(criteria); err != nil {
		return nil, fmt.Errorf("invalid search criteria: %w", err)
	}

	// Build query
	query := `SELECT id, service, operation, start_time, end_time, duration, status, tags FROM traces WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if criteria.Service != nil {
		query += fmt.Sprintf(" AND service = $%d", argIndex)
		args = append(args, *criteria.Service)
		argIndex++
	}

	if criteria.Operation != nil {
		query += fmt.Sprintf(" AND operation = $%d", argIndex)
		args = append(args, *criteria.Operation)
		argIndex++
	}

	if criteria.StartTime != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIndex)
		args = append(args, *criteria.StartTime)
		argIndex++
	}

	if criteria.EndTime != nil {
		query += fmt.Sprintf(" AND start_time <= $%d", argIndex)
		args = append(args, *criteria.EndTime)
		argIndex++
	}

	if criteria.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *criteria.Status)
		argIndex++
	}

	// Add ordering and pagination
	query += " ORDER BY start_time DESC"
	
	if criteria.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, criteria.Limit)
		argIndex++
	}

	if criteria.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, criteria.Offset)
	}

	// Execute query
	rows, err := tr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query traces: %w", err)
	}
	defer rows.Close()

	var traces []*domain.Trace
	for rows.Next() {
		var trace domain.Trace
		var tagsJSON []byte

		err := rows.Scan(
			&trace.ID,
			&trace.Service,
			&trace.Operation,
			&trace.StartTime,
			&trace.EndTime,
			&trace.Duration,
			&trace.Status,
			&tagsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trace: %w", err)
		}

		// Parse tags
		if err := json.Unmarshal(tagsJSON, &trace.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		traces = append(traces, &trace)
	}

	return traces, nil
}

// GetServices returns all available services
func (tr *traceRepositoryPostgres) GetServices(ctx context.Context) ([]domain.ServiceName, error) {
	query := `SELECT DISTINCT service FROM traces ORDER BY service`
	
	rows, err := tr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []domain.ServiceName
	for rows.Next() {
		var service string
		if err := rows.Scan(&service); err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, domain.ServiceName(service))
	}

	return services, nil
}

// GetOperations returns all operations for a specific service
func (tr *traceRepositoryPostgres) GetOperations(ctx context.Context, service domain.ServiceName) ([]domain.OperationName, error) {
	if service == "" {
		return nil, fmt.Errorf("service name is required")
	}

	query := `SELECT DISTINCT operation FROM traces WHERE service = $1 ORDER BY operation`
	
	rows, err := tr.db.QueryContext(ctx, query, service)
	if err != nil {
		return nil, fmt.Errorf("failed to query operations: %w", err)
	}
	defer rows.Close()

	var operations []domain.OperationName
	for rows.Next() {
		var operation string
		if err := rows.Scan(&operation); err != nil {
			return nil, fmt.Errorf("failed to scan operation: %w", err)
		}
		operations = append(operations, domain.OperationName(operation))
	}

	return operations, nil
}

// validateTrace validates a trace before saving
func (tr *traceRepositoryPostgres) validateTrace(trace *domain.Trace) error {
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
func (tr *traceRepositoryPostgres) validateSearchCriteria(criteria *domain.SearchCriteria) error {
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

// Close closes the database connection
func (tr *traceRepositoryPostgres) Close() error {
	return tr.db.Close()
}
