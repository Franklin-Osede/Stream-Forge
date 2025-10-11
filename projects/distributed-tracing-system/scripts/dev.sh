#!/bin/bash

# Development script for distributed tracing system
# Uses Docker to run tests and build without requiring Go installation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run tests
run_tests() {
    print_status "Running unit tests..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "apk add --no-cache gcc musl-dev && go mod tidy && CGO_ENABLED=1 go test -v -race -short ./internal/..."
    
    print_status "Running integration tests..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "apk add --no-cache gcc musl-dev && go mod tidy && CGO_ENABLED=1 go test -v -tags=integration ./tests/integration/..."
}

# Function to run tests with coverage
run_tests_coverage() {
    print_status "Running tests with coverage..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "go test -v -race -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"
    
    print_success "Coverage report generated: coverage.html"
}

# Function to build the application
build_app() {
    print_status "Building application..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "go build -o bin/distributed-tracing-system ./cmd/server"
    
    print_success "Application built successfully"
}

# Function to run linter
run_linter() {
    print_status "Running linter..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && golangci-lint run"
}

# Function to format code
format_code() {
    print_status "Formatting code..."
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
        sh -c "go fmt ./... && go install golang.org/x/tools/cmd/goimports@latest && goimports -w ."
    
    print_success "Code formatted successfully"
}

# Function to run TDD workflow
run_tdd() {
    print_status "Starting TDD workflow..."
    print_status "1. Running tests (RED phase)..."
    run_tests
    
    print_status "2. Building application (GREEN phase)..."
    build_app
    
    print_status "3. Running linter (BLUE phase)..."
    run_linter
    
    print_success "TDD workflow completed!"
}

# Function to start development environment
start_dev() {
    print_status "Starting development environment..."
    
    # Start Docker services
    print_status "Starting Docker services..."
    docker-compose up -d
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check service health
    print_status "Checking service health..."
    curl -s http://localhost:16686 > /dev/null && print_success "Jaeger is running" || print_warning "Jaeger not ready"
    curl -s http://localhost:9090 > /dev/null && print_success "Prometheus is running" || print_warning "Prometheus not ready"
    curl -s http://localhost:8080 > /dev/null && print_success "Kafka UI is running" || print_warning "Kafka UI not ready"
    
    print_success "Development environment started!"
    print_status "Services available at:"
    print_status "  - Jaeger UI: http://localhost:16686"
    print_status "  - Prometheus: http://localhost:9090"
    print_status "  - Kafka UI: http://localhost:8080"
    print_status "  - Grafana: http://localhost:3000"
}

# Function to stop development environment
stop_dev() {
    print_status "Stopping development environment..."
    docker-compose down
    print_success "Development environment stopped!"
}

# Function to show help
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  test        Run unit tests"
    echo "  test-integration  Run integration tests"
    echo "  test-coverage    Run tests with coverage"
    echo "  build       Build the application"
    echo "  lint        Run linter"
    echo "  format      Format code"
    echo "  tdd         Run TDD workflow"
    echo "  dev         Start development environment"
    echo "  dev-stop    Stop development environment"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 test     # Run unit tests"
    echo "  $0 tdd      # Run TDD workflow"
    echo "  $0 dev      # Start development environment"
}

# Main script logic
case "${1:-help}" in
    test)
        run_tests
        ;;
    test-integration)
        docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine \
            sh -c "go test -v -tags=integration ./tests/integration/..."
        ;;
    test-coverage)
        run_tests_coverage
        ;;
    build)
        build_app
        ;;
    lint)
        run_linter
        ;;
    format)
        format_code
        ;;
    tdd)
        run_tdd
        ;;
    dev)
        start_dev
        ;;
    dev-stop)
        stop_dev
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
