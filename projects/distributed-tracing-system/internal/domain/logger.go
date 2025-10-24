package domain

import (
	"context"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// Logger defines the interface for structured logging
type Logger interface {
	// Basic logging methods
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// Context-aware logging
	WithContext(ctx context.Context) Logger
	WithFields(fields ...Field) Logger

	// Structured logging with levels
	Log(level LogLevel, msg string, fields ...Field)

	// Utility methods
	SetLevel(level LogLevel)
	GetLevel() LogLevel
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new field
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level      LogLevel `json:"level"`
	Format     string   `json:"format"`     // json, text
	Output     string   `json:"output"`     // stdout, stderr, file
	FilePath   string   `json:"file_path"`  // for file output
	MaxSize    int      `json:"max_size"`   // in MB
	MaxBackups int      `json:"max_backups"`
	MaxAge     int      `json:"max_age"` // in days
	Compress   bool     `json:"compress"`
}

// GetDefaultConfig returns default logger configuration
func GetDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      InfoLevel,
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
}
