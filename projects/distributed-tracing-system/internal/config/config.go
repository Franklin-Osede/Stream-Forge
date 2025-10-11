package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Jaeger   JaegerConfig
	Kafka    KafkaConfig
	Prometheus PrometheusConfig
	Logging  LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// JaegerConfig holds Jaeger configuration
type JaegerConfig struct {
	Endpoint string
	Timeout  time.Duration
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers       []string
	TopicTraces   string
	GroupID       string
	RetryAttempts int
	RetryDelay    time.Duration
}

// PrometheusConfig holds Prometheus configuration
type PrometheusConfig struct {
	Port string
	Path string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Jaeger: JaegerConfig{
			Endpoint: getEnv("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
			Timeout:  getDurationEnv("JAEGER_TIMEOUT", 30*time.Second),
		},
		Kafka: KafkaConfig{
			Brokers:       getStringSliceEnv("KAFKA_BROKERS", []string{"localhost:9092"}),
			TopicTraces:   getEnv("KAFKA_TOPIC_TRACES", "trace-events"),
			GroupID:       getEnv("KAFKA_GROUP_ID", "tracing-system"),
			RetryAttempts: getIntEnv("KAFKA_RETRY_ATTEMPTS", 3),
			RetryDelay:    getDurationEnv("KAFKA_RETRY_DELAY", 1*time.Second),
		},
		Prometheus: PrometheusConfig{
			Port: getEnv("PROMETHEUS_PORT", "9091"),
			Path: getEnv("PROMETHEUS_PATH", "/metrics"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return cfg, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values
		// In production, you might want to use a more sophisticated parser
		return []string{value}
	}
	return defaultValue
}
