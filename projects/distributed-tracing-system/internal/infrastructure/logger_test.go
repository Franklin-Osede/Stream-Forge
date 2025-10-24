package infrastructure

import (
	"context"
	"testing"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZapLogger_Creation(t *testing.T) {
	tests := []struct {
		name   string
		config *domain.LoggerConfig
		wantErr bool
	}{
		{
			name:     "default config",
			config:   nil,
			wantErr:  false,
		},
		{
			name: "custom config",
			config: &domain.LoggerConfig{
				Level:  domain.InfoLevel,
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "debug level",
			config: &domain.LoggerConfig{
				Level:  domain.DebugLevel,
				Format: "text",
				Output: "stdout",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewZapLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}

func TestZapLogger_LogLevels(t *testing.T) {
	config := &domain.LoggerConfig{
		Level:  domain.DebugLevel,
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// Test all log levels
	logger.Debug("debug message", domain.NewField("test", "debug"))
	logger.Info("info message", domain.NewField("test", "info"))
	logger.Warn("warn message", domain.NewField("test", "warn"))
	logger.Error("error message", domain.NewField("test", "error"))
	
	// Fatal would exit the process, so we skip it in tests
}

func TestZapLogger_WithFields(t *testing.T) {
	config := &domain.LoggerConfig{
		Level:  domain.InfoLevel,
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// Test with fields
	loggerWithFields := logger.WithFields(
		domain.NewField("service", "test-service"),
		domain.NewField("version", "1.0.0"),
	)

	assert.NotNil(t, loggerWithFields)
	
	// Test that the new logger has the fields
	loggerWithFields.Info("message with fields")
}

func TestZapLogger_WithContext(t *testing.T) {
	config := &domain.LoggerConfig{
		Level:  domain.InfoLevel,
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// Test with context
	ctx := context.WithValue(context.Background(), "correlation_id", "test-correlation-123")
	ctx = context.WithValue(ctx, "request_id", "test-request-456")

	loggerWithContext := logger.WithContext(ctx)
	assert.NotNil(t, loggerWithContext)
	
	// Test that the new logger has context information
	loggerWithContext.Info("message with context")
}

func TestZapLogger_LogMethod(t *testing.T) {
	config := &domain.LoggerConfig{
		Level:  domain.InfoLevel,
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// Test Log method with different levels
	logger.Log(domain.DebugLevel, "debug via Log method", domain.NewField("level", "debug"))
	logger.Log(domain.InfoLevel, "info via Log method", domain.NewField("level", "info"))
	logger.Log(domain.WarnLevel, "warn via Log method", domain.NewField("level", "warn"))
	logger.Log(domain.ErrorLevel, "error via Log method", domain.NewField("level", "error"))
}

func TestZapLogger_SetLevel(t *testing.T) {
	config := &domain.LoggerConfig{
		Level:  domain.InfoLevel,
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// Test initial level
	assert.Equal(t, domain.InfoLevel, logger.GetLevel())

	// Test setting level
	logger.SetLevel(domain.DebugLevel)
	assert.Equal(t, domain.DebugLevel, logger.GetLevel())

	logger.SetLevel(domain.ErrorLevel)
	assert.Equal(t, domain.ErrorLevel, logger.GetLevel())
}

func TestLoggerFactory_CreateLogger(t *testing.T) {
	factory := NewLoggerFactory()

	tests := []struct {
		name   string
		config *domain.LoggerConfig
		wantErr bool
	}{
		{
			name:     "default config",
			config:   nil,
			wantErr:  false,
		},
		{
			name: "custom config",
			config: &domain.LoggerConfig{
				Level:  domain.InfoLevel,
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := factory.CreateLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}

func TestLoggerFactory_CreateDefaultLogger(t *testing.T) {
	factory := NewLoggerFactory()

	logger, err := factory.CreateDefaultLogger()
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestLoggerFactory_CreateLoggerForService(t *testing.T) {
	factory := NewLoggerFactory()

	logger, err := factory.CreateLoggerForService("test-service", domain.InfoLevel)
	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Test that the logger has service information
	logger.Info("service message")
}

func TestLoggerFactory_CreateLoggerForComponent(t *testing.T) {
	factory := NewLoggerFactory()

	logger, err := factory.CreateLoggerForComponent("test-component", domain.InfoLevel)
	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Test that the logger has component information
	logger.Info("component message")
}
