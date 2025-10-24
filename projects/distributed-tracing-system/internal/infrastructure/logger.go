package infrastructure

import (
	"context"
	"fmt"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
	level  domain.LogLevel
}

// NewZapLogger creates a new zap-based logger
func NewZapLogger(config *domain.LoggerConfig) (domain.Logger, error) {
	if config == nil {
		config = domain.GetDefaultLoggerConfig()
	}

	// Create zap configuration
	zapConfig := zap.NewProductionConfig()
	
	// Set log level
	zapConfig.Level = zap.NewAtomicLevelAt(domainToZapLevel(config.Level))
	
	// Set output
	switch config.Output {
	case "stderr":
		zapConfig.OutputPaths = []string{"stderr"}
	case "file":
		if config.FilePath == "" {
			config.FilePath = "logs/app.log"
		}
		zapConfig.OutputPaths = []string{config.FilePath}
	default:
		zapConfig.OutputPaths = []string{"stdout"}
	}

	// Set format
	if config.Format == "text" {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Create logger
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return &zapLogger{
		logger: logger,
		level:  config.Level,
	}, nil
}

// Debug logs a debug message
func (z *zapLogger) Debug(msg string, fields ...domain.Field) {
	if z.level <= domain.DebugLevel {
		z.logger.Debug(msg, z.convertFields(fields)...)
	}
}

// Info logs an info message
func (z *zapLogger) Info(msg string, fields ...domain.Field) {
	if z.level <= domain.InfoLevel {
		z.logger.Info(msg, z.convertFields(fields)...)
	}
}

// Warn logs a warning message
func (z *zapLogger) Warn(msg string, fields ...domain.Field) {
	if z.level <= domain.WarnLevel {
		z.logger.Warn(msg, z.convertFields(fields)...)
	}
}

// Error logs an error message
func (z *zapLogger) Error(msg string, fields ...domain.Field) {
	if z.level <= domain.ErrorLevel {
		z.logger.Error(msg, z.convertFields(fields)...)
	}
}

// Fatal logs a fatal message and exits
func (z *zapLogger) Fatal(msg string, fields ...domain.Field) {
	z.logger.Fatal(msg, z.convertFields(fields)...)
}

// WithContext creates a logger with context information
func (z *zapLogger) WithContext(ctx context.Context) domain.Logger {
	fields := []domain.Field{}

	// Add trace context if available
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		if spanCtx.IsValid() {
			fields = append(fields,
				domain.NewField("trace_id", spanCtx.TraceID().String()),
				domain.NewField("span_id", spanCtx.SpanID().String()),
			)
		}
	}

	// Add correlation ID if available
	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		fields = append(fields, domain.NewField("correlation_id", correlationID))
	}

	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, domain.NewField("request_id", requestID))
	}

	return z.WithFields(fields...)
}

// WithFields creates a logger with additional fields
func (z *zapLogger) WithFields(fields ...domain.Field) domain.Logger {
	zapFields := z.convertFields(fields)
	newLogger := z.logger.With(zapFields...)
	
	return &zapLogger{
		logger: newLogger,
		level:  z.level,
	}
}

// Log logs a message with the specified level
func (z *zapLogger) Log(level domain.LogLevel, msg string, fields ...domain.Field) {
	switch level {
	case domain.DebugLevel:
		z.Debug(msg, fields...)
	case domain.InfoLevel:
		z.Info(msg, fields...)
	case domain.WarnLevel:
		z.Warn(msg, fields...)
	case domain.ErrorLevel:
		z.Error(msg, fields...)
	case domain.FatalLevel:
		z.Fatal(msg, fields...)
	}
}

// SetLevel sets the log level
func (z *zapLogger) SetLevel(level domain.LogLevel) {
	z.level = level
}

// GetLevel returns the current log level
func (z *zapLogger) GetLevel() domain.LogLevel {
	return z.level
}

// convertFields converts domain fields to zap fields
func (z *zapLogger) convertFields(fields []domain.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// domainToZapLevel converts domain log level to zap level
func domainToZapLevel(level domain.LogLevel) zapcore.Level {
	switch level {
	case domain.DebugLevel:
		return zapcore.DebugLevel
	case domain.InfoLevel:
		return zapcore.InfoLevel
	case domain.WarnLevel:
		return zapcore.WarnLevel
	case domain.ErrorLevel:
		return zapcore.ErrorLevel
	case domain.FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Sync flushes any buffered log entries
func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}

// Close closes the logger
func (z *zapLogger) Close() error {
	return z.logger.Sync()
}
