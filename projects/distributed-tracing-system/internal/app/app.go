package app

import (
	"context"
	"fmt"
	"log"

	"github.com/streamforge/distributed-tracing-system/internal/config"
	"github.com/streamforge/distributed-tracing-system/internal/infrastructure"
	"github.com/streamforge/distributed-tracing-system/internal/interfaces"
	"github.com/streamforge/distributed-tracing-system/internal/usecases"
)

// App represents the application
type App struct {
	config *config.Config
	server *interfaces.Server
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize infrastructure
	jaegerExporter, err := infrastructure.NewJaegerExporter(cfg.Jaeger.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	prometheusExporter, err := infrastructure.NewPrometheusExporter(cfg.Prometheus.Port, cfg.Prometheus.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	kafkaProducer, err := infrastructure.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.TopicTraces)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	kafkaConsumer, err := infrastructure.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.TopicTraces, cfg.Kafka.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	// Initialize repositories
	traceRepo := infrastructure.NewTraceRepository(jaegerExporter)

	// Initialize use cases
	traceService := usecases.NewTraceService(traceRepo, prometheusExporter, kafkaProducer)

	// Initialize interfaces
	server, err := interfaces.NewServer(cfg, traceService)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	// Start Kafka consumer
	go func() {
		if err := kafkaConsumer.Start(context.Background(), traceService); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	return &App{
		config: cfg,
		server: server,
	}, nil
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	log.Printf("Starting distributed tracing system on port %s", a.config.Server.Port)
	
	return a.server.Start(ctx)
}

