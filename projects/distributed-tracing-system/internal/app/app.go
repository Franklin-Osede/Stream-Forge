package app

import (
	"context"
	"fmt"

	"github.com/streamforge/distributed-tracing-system/internal/config"
	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/streamforge/distributed-tracing-system/internal/infrastructure"
	"github.com/streamforge/distributed-tracing-system/internal/interfaces"
	"github.com/streamforge/distributed-tracing-system/internal/telemetry"
	"github.com/streamforge/distributed-tracing-system/internal/usecases"
)

// App represents the application
type App struct {
	config *config.Config
	server *interfaces.ServerWithTelemetry
	logger domain.Logger
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	loggerFactory := infrastructure.NewLoggerFactory()
	logger, err := loggerFactory.CreateLoggerForService("distributed-tracing-system", domain.InfoLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("Starting distributed tracing system", 
		domain.NewField("version", "1.0.0"),
		domain.NewField("environment", "development"),
	)

	// Initialize telemetry
	telemetryConfig := &telemetry.TelemetryConfig{
		ServiceName:    "distributed-tracing-system",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		JaegerEndpoint: cfg.Jaeger.Endpoint,
		PrometheusPort: cfg.Prometheus.Port,
		SamplingRate:   0.1, // 10% sampling
	}

	telemetryManager, err := telemetry.NewTelemetryManager(telemetryConfig)
	if err != nil {
		logger.Error("Failed to create telemetry manager", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create telemetry manager: %w", err)
	}

	logger.Info("Telemetry manager initialized successfully")

	// Initialize Jaeger exporter with complete OpenTelemetry SDK
	jaegerExporter, err := infrastructure.NewJaegerExporterComplete(telemetryManager)
	if err != nil {
		logger.Error("Failed to create Jaeger exporter", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	logger.Info("Jaeger exporter initialized successfully")

	prometheusExporter, err := infrastructure.NewPrometheusExporter(cfg.Prometheus.Port, cfg.Prometheus.Path)
	if err != nil {
		logger.Error("Failed to create Prometheus exporter", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	logger.Info("Prometheus exporter initialized successfully")

	kafkaProducer, err := infrastructure.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.TopicTraces)
	if err != nil {
		logger.Error("Failed to create Kafka producer", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	logger.Info("Kafka producer initialized successfully")

	kafkaConsumer, err := infrastructure.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.TopicTraces, cfg.Kafka.GroupID)
	if err != nil {
		logger.Error("Failed to create Kafka consumer", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	logger.Info("Kafka consumer initialized successfully")

	// Initialize repositories
	traceRepo, err := infrastructure.NewTraceRepositoryPostgres(cfg.Database.GetDSN(), jaegerExporter)
	if err != nil {
		logger.Error("Failed to create trace repository", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create trace repository: %w", err)
	}

	logger.Info("Trace repository initialized successfully")

	// Initialize use cases
	traceService := usecases.NewTraceService(traceRepo, prometheusExporter, kafkaProducer)

	// Initialize interfaces with telemetry
	server, err := interfaces.NewServerWithTelemetry(cfg, traceService, telemetryManager)
	if err != nil {
		logger.Error("Failed to create server", domain.NewField("error", err.Error()))
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	logger.Info("Server initialized successfully")

	// Start Kafka consumer
	go func() {
		if err := kafkaConsumer.Start(context.Background(), traceService); err != nil {
			logger.Error("Kafka consumer error", domain.NewField("error", err.Error()))
		}
	}()

	logger.Info("Application initialized successfully")

	return &App{
		config: cfg,
		server: server,
		logger: logger,
	}, nil
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting distributed tracing system", 
		domain.NewField("port", a.config.Server.Port),
		domain.NewField("environment", "development"),
	)
	
	return a.server.Start(ctx)
}

