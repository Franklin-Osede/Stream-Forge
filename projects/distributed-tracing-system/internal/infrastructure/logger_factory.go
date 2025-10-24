package infrastructure

import (
	"fmt"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
)

// LoggerFactory creates logger instances
type LoggerFactory struct{}

// NewLoggerFactory creates a new logger factory
func NewLoggerFactory() *LoggerFactory {
	return &LoggerFactory{}
}

// CreateLogger creates a logger based on the configuration
func (lf *LoggerFactory) CreateLogger(config *domain.LoggerConfig) (domain.Logger, error) {
	if config == nil {
		config = domain.GetDefaultLoggerConfig()
	}

	// For now, we only support zap logger
	// In the future, we could add support for other loggers (logrus, slog, etc.)
	return NewZapLogger(config)
}

// CreateDefaultLogger creates a logger with default configuration
func (lf *LoggerFactory) CreateDefaultLogger() (domain.Logger, error) {
	return lf.CreateLogger(nil)
}

// CreateLoggerForService creates a logger configured for a specific service
func (lf *LoggerFactory) CreateLoggerForService(serviceName string, level domain.LogLevel) (domain.Logger, error) {
	config := &domain.LoggerConfig{
		Level:  level,
		Format: "json",
		Output: "stdout",
	}

	logger, err := lf.CreateLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger for service %s: %w", serviceName, err)
	}

	// Add service name to all logs
	return logger.WithFields(domain.NewField("service", serviceName)), nil
}

// CreateLoggerForComponent creates a logger for a specific component
func (lf *LoggerFactory) CreateLoggerForComponent(component string, level domain.LogLevel) (domain.Logger, error) {
	config := &domain.LoggerConfig{
		Level:  level,
		Format: "json",
		Output: "stdout",
	}

	logger, err := lf.CreateLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger for component %s: %w", component, err)
	}

	// Add component name to all logs
	return logger.WithFields(domain.NewField("component", component)), nil
}
