package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streamforge/distributed-tracing-system/internal/app"
	"github.com/streamforge/distributed-tracing-system/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Run application
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

