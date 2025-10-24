# Distributed Tracing System Makefile

.PHONY: help test test-unit test-integration test-coverage test-watch build run clean lint format

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Testing targets
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	go test -v -race -short ./internal/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-watch: ## Run tests in watch mode
	@echo "Running tests in watch mode..."
	@which entr > /dev/null || (echo "Please install entr: brew install entr" && exit 1)
	find . -name "*.go" | entr -c make test-unit

# Build targets
build: ## Build the application
	@echo "Building distributed tracing system..."
	go build -o bin/distributed-tracing-system ./cmd/server

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/distributed-tracing-system-linux ./cmd/server

build-docker: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t streamforge/distributed-tracing-system:latest .

# Run targets
run: build ## Build and run the application
	@echo "Running distributed tracing system..."
	./bin/distributed-tracing-system

run-dev: ## Run in development mode with hot reload
	@echo "Running in development mode..."
	@which air > /dev/null || (echo "Please install air: go install github.com/cosmtrek/air@latest" && exit 1)
	air -c .air.conf

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Please install golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

clean-docker: ## Clean Docker artifacts
	@echo "Cleaning Docker artifacts..."
	docker rmi streamforge/distributed-tracing-system:latest || true

# Development targets
deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# TDD specific targets
tdd: test-watch ## Run TDD workflow (test-driven development)

tdd-setup: ## Setup TDD environment
	@echo "Setting up TDD environment..."
	go install github.com/cosmtrek/air@latest
	go install golangci-lint/cmd/golangci-lint@latest
	@which entr > /dev/null || (echo "Please install entr: brew install entr")

# Docker targets
docker-up: ## Start Docker services for development
	@echo "Starting Docker services..."
	docker-compose -f docker-compose.dev.yml up -d postgres jaeger prometheus kafka

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker-compose -f docker-compose.dev.yml down

docker-logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	docker-compose -f docker-compose.dev.yml logs -f

# Database targets
db-up: ## Start PostgreSQL database
	@echo "Starting PostgreSQL database..."
	docker-compose -f docker-compose.dev.yml up -d postgres

db-down: ## Stop PostgreSQL database
	@echo "Stopping PostgreSQL database..."
	docker-compose -f docker-compose.dev.yml stop postgres

db-logs: ## Show database logs
	@echo "Showing database logs..."
	docker-compose -f docker-compose.dev.yml logs -f postgres

db-connect: ## Connect to PostgreSQL database
	@echo "Connecting to PostgreSQL database..."
	docker exec -it tracing-postgres psql -U postgres -d tracing_system

# Health checks
health: ## Check service health
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health || echo "Service not running"
	@curl -s http://localhost:16686 || echo "Jaeger not running"
	@curl -s http://localhost:9090 || echo "Prometheus not running"

# Development workflow
dev: docker-up run-dev ## Start development environment

dev-stop: docker-down ## Stop development environment

# CI/CD targets
ci: test lint build ## Run CI pipeline

ci-docker: build-docker ## Build Docker image for CI

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	@which godoc > /dev/null || (echo "Please install godoc: go install golang.org/x/tools/cmd/godoc@latest" && exit 1)
	godoc -http=:6060

# Performance testing
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Security scanning
security: ## Run security scan
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Please install gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

# All-in-one targets
all: clean deps test lint build ## Run everything (clean, deps, test, lint, build)

all-docker: clean deps test lint build-docker ## Run everything with Docker build

