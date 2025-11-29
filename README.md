<div align="center">

# StreamForge

**Enterprise-grade microservices ecosystem for stream processing, real-time data analytics, and distributed observability**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://golang.org/)
[![Prometheus](https://img.shields.io/badge/Monitoring-Prometheus-red.svg)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/Dashboard-Grafana-orange.svg)](https://grafana.com/)
[![Jaeger](https://img.shields.io/badge/Tracing-Jaeger-blue.svg)](https://www.jaegertracing.io/)

[Documentation](#documentation) • [Quick Start](#getting-started) • [Architecture](#architecture) • [Projects](#projects) • [Contributing](#contributing)

</div>

---

## Project Status

**Note: This project is currently under active development. Not all features are complete, and some components may be in various stages of implementation. Please refer to the individual project READMEs for current status and known limitations.**

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Projects](#projects)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Development](#development)
- [Docker & Deployment](#docker--deployment)
- [Monitoring & Observability](#monitoring--observability)
- [Security](#security)
- [Documentation](#documentation)
- [Testing](#testing)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Overview

StreamForge is a production-ready microservices ecosystem designed for enterprise-grade stream processing, real-time data analytics, and comprehensive observability. The platform provides a complete solution for processing high-volume data streams, detecting anomalies, building analytics dashboards, and maintaining full visibility across distributed systems.

**Core Capabilities**

- **Real-time Stream Processing** - Kafka-based event streaming with KSQLDB and Flink
- **Distributed Tracing** - OpenTelemetry integration with Jaeger for end-to-end request tracking
- **Machine Learning** - TensorFlow-powered anomaly detection in real-time
- **Interactive Dashboards** - Drag-and-drop dashboard builder with Angular
- **Intelligent Alerting** - ML-driven alert management to reduce alert fatigue
- **Multi-tenant Architecture** - Complete isolation between tenants
- **Full Observability** - Prometheus metrics, Grafana dashboards, and Jaeger tracing

## Key Features

### Stream Processing
- Kafka-based event streaming architecture
- KSQLDB for stream processing queries
- Apache Flink for complex event processing
- Real-time data validation and transformation
- Schema registry for data governance

### Observability
- Distributed tracing with OpenTelemetry
- Prometheus metrics collection
- Grafana dashboards and visualizations
- Jaeger trace visualization
- Structured logging with correlation IDs
- Health checks and readiness probes

### Machine Learning
- Real-time anomaly detection
- TensorFlow model integration
- Isolation Forest and LSTM models
- Model monitoring and drift detection

### User Interface
- Real-time dashboard builder (drag-and-drop)
- Metrics portal with Grafana integration
- WebSocket and Server-Sent Events for live updates
- Export capabilities (PDF, PNG, JSON)
- Pre-built dashboard templates

### Enterprise Features
- Multi-tenant isolation engine
- Data validation and quality checks
- Audit logging and compliance
- Rate limiting and circuit breakers
- Security best practices

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Data Sources Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Blockchain   │  │ AI Models    │  │ IoT Sensors │        │
│  │ APIs         │  │              │  │             │        │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘        │
└─────────┼──────────────────┼──────────────────┼───────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│              Event Processing Layer                             │
│  ┌──────────────────────┐  ┌──────────────────────┐            │
│  │ event-bridge-kafka   │  │ stream-data-validator│            │
│  │ (NestJS)             │  │ (Go)                 │            │
│  └──────────┬───────────┘  └──────────┬───────────┘            │
└─────────────┼──────────────────────────┼─────────────────────────┘
              │                        │
              ▼                        ▼
┌─────────────────────────────────────────────────────────────────┐
│              Stream Processing Layer                            │
│  ┌──────────────────────┐  ┌──────────────────────┐            │
│  │ stream-analytics-hub │  │ stream-anomaly-     │            │
│  │ (KSQLDB + Flink)     │  │ detector (FastAPI)  │            │
│  └──────────┬───────────┘  └──────────┬───────────┘            │
└─────────────┼──────────────────────────┼─────────────────────────┘
              │                        │
              ▼                        ▼
┌─────────────────────────────────────────────────────────────────┐
│              Observability Layer                                 │
│  ┌──────────────────────┐  ┌──────────────────────┐            │
│  │ distributed-tracing- │  │ intelligent-alert-   │            │
│  │ system (Go)          │  │ manager (Python)    │            │
│  └──────────┬───────────┘  └──────────┬───────────┘            │
└─────────────┼──────────────────────────┼─────────────────────────┘
              │                        │
              ▼                        ▼
┌─────────────────────────────────────────────────────────────────┐
│              Frontend Layer                                      │
│  ┌──────────────────────┐  ┌──────────────────────┐            │
│  │ kafka-metrics-portal  │  │ real-time-dashboard- │            │
│  │ (React + NestJS)      │  │ builder (Angular)   │            │
│  └──────────────────────┘  └──────────────────────┘            │
└─────────────────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────────┐
│              Infrastructure Layer                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ Kafka    │  │PostgreSQL│  │ Redis    │  │ Jaeger   │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                      │
│  │Prometheus│  │ Grafana  │  │ KSQLDB   │                      │
│  └──────────┘  └──────────┘  └──────────┘                      │
└─────────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Backend Languages:**
- **Go** - distributed-tracing-system, stream-data-validator, multi-tenant-isolation-engine
- **Python** - stream-anomaly-detector, intelligent-alert-manager, log-replay-simulator
- **Node.js/TypeScript** - event-bridge-kafka, kafka-metrics-portal
- **Java** - stream-analytics-hub (KSQLDB + Flink)

**Frameworks & Libraries:**
- **NestJS** - REST APIs and microservices
- **FastAPI** - High-performance Python APIs
- **Gin** - Lightweight Go web framework
- **Spring Boot** - Enterprise Java applications
- **OpenTelemetry** - Distributed tracing instrumentation

**Frontend:**
- **Angular 17+** - Dashboard builder
- **React** - Metrics portal
- **TypeScript** - Type-safe development
- **D3.js, Chart.js** - Data visualization

**Infrastructure:**
- **Docker & Docker Compose** - Containerization
- **Apache Kafka** - Event streaming platform
- **KSQLDB** - Stream processing SQL engine
- **Apache Flink** - Complex event processing
- **PostgreSQL** - Relational database
- **Redis** - Caching and session storage

**Observability:**
- **Prometheus** - Metrics collection
- **Grafana** - Visualization and dashboards
- **Jaeger** - Distributed tracing
- **OpenTelemetry** - Instrumentation standard

## Projects

### Core Projects

| # | Project | Description | Stack | Status |
|---|---------|-------------|-------|--------|
| 1 | **[event-bridge-kafka](./projects/event-bridge-kafka/)** | Blockchain/AI events gateway and router | NestJS • KafkaJS • Docker | Complete |
| 2 | **[stream-anomaly-detector](./projects/stream-anomaly-detector/)** | Real-time ML anomaly detection | FastAPI • TensorFlow • Kafka | Complete |
| 3 | **[stream-analytics-hub](./projects/stream-analytics-hub/)** | KSQLDB + Flink analytics engine | KSQLDB • Flink • Prometheus | Complete |
| 4 | **[kafka-metrics-portal](./projects/kafka-metrics-portal/)** | Metrics visualization portal | React • NestJS • Grafana | Complete |
| 5 | **[log-replay-simulator](./projects/log-replay-simulator/)** | Traffic and event simulator | Python • Kafka • Docker | Complete |

### Advanced Projects

| # | Project | Description | Stack | Status |
|---|---------|-------------|-------|--------|
| 6 | **[distributed-tracing-system](./projects/distributed-tracing-system/)** | Distributed tracing with OpenTelemetry | Go • OpenTelemetry • Jaeger • PostgreSQL | In Progress |
| 7 | **[intelligent-alert-manager](./projects/intelligent-alert-manager/)** | ML-powered alert management | Python • TensorFlow • Redis | Planned |
| 8 | **[real-time-dashboard-builder](./projects/real-time-dashboard-builder/)** | Drag-and-drop dashboard builder | Angular • D3.js • WebSocket | Planned |
| 9 | **[stream-data-validator](./projects/stream-data-validator/)** | Real-time data validation | Go • Avro • Kafka | Planned |
| 10 | **[multi-tenant-isolation-engine](./projects/multi-tenant-isolation-engine/)** | Multi-tenant isolation and security | Go • Kubernetes • Istio • Vault | Planned |

**Note:** Project statuses reflect current development state. Some projects marked as "Complete" may still have ongoing improvements, and "In Progress" projects are actively being developed. Check individual project READMEs for detailed status.

## Getting Started

### Prerequisites

- **Docker** & **Docker Compose** (v2.0+)
- **Git** (for cloning the repository)
- **Make** (optional, for convenience commands)
- **Go 1.23+** (for distributed-tracing-system development)
- **Node.js 18+** (for Node.js-based projects)
- **Python 3.10+** (for Python-based projects)

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/your-org/stream-forge.git
cd stream-forge

# 2. Start the entire ecosystem
make up

# 3. Check service status
make status

# 4. View logs
make logs
```

### Service Access

Once the services are running, access them at:

| Service | URL | Description | Credentials |
|---------|-----|-------------|-------------|
| **Grafana** | http://localhost:3000 | Dashboards and visualizations | admin / admin123 |
| **Prometheus** | http://localhost:9090 | Metrics and alerts | - |
| **Jaeger** | http://localhost:16686 | Distributed tracing UI | - |
| **Kafka UI** | http://localhost:8080 | Kafka management (if enabled) | - |
| **distributed-tracing-system** | http://localhost:8082 | Tracing API | - |

### Useful Commands

```bash
# Show all available commands
make help

# Start specific project
make up-project PROJECT=distributed-tracing-system

# View logs for specific project
make logs-project PROJECT=distributed-tracing-system

# Stop all services
make down

# Clean everything (containers, volumes, networks)
make clean

# Build all Docker images
make build

# Run tests across all projects
make test

# Format code
make format

# Lint code
make lint

# Open monitoring dashboards
make monitor
```

## Development

### Distributed Tracing System (Go)

The distributed tracing system is the core observability component, built with Go and OpenTelemetry. This project is currently in active development.

#### Project Structure

```
distributed-tracing-system/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/                     # Application initialization
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain entities and interfaces
│   ├── infrastructure/          # External integrations
│   │   ├── jaeger_exporter.go   # Jaeger exporter
│   │   ├── kafka_consumer.go    # Kafka consumer
│   │   ├── kafka_producer.go    # Kafka producer
│   │   ├── logger.go            # Structured logging
│   │   └── prometheus_exporter.go # Prometheus metrics
│   ├── interfaces/              # HTTP handlers
│   ├── telemetry/               # Telemetry middleware
│   └── usecases/                # Business logic
├── tests/
│   └── integration/            # Integration tests
├── Dockerfile                   # Container image
├── docker-compose.yml           # Local development
└── go.mod                       # Go dependencies
```

#### Key Features

- OpenTelemetry SDK integration
- Jaeger exporter for trace visualization
- Prometheus metrics exporter
- Kafka consumer for trace events
- PostgreSQL repository for trace storage
- Structured logging with correlation IDs
- Health check endpoints
- RESTful API for trace queries

**Note:** Some features may be partially implemented or under active development. Refer to the project's README for current implementation status.

#### API Endpoints

```yaml
GET  /health                          # Health check
GET  /api/v1/traces/search            # Search traces
GET  /api/v1/traces/:id               # Get trace by ID
GET  /api/v1/services                 # List all services
GET  /api/v1/services/:service/operations # List operations for service
GET  /api/v1/metrics                  # Get tracing metrics
```

#### Configuration

Environment variables:

```bash
# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=streamforge
DB_PASSWORD=streamforge123
DB_NAME=tracing_system
DB_SSLMODE=disable

# Jaeger
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
JAEGER_TIMEOUT=30s

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_TRACES=trace-events
KAFKA_GROUP_ID=tracing-system

# Prometheus
PROMETHEUS_PORT=9091
PROMETHEUS_PATH=/metrics

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

#### Local Development

```bash
cd projects/distributed-tracing-system

# Install Go dependencies
go mod download

# Run locally
go run cmd/server/main.go

# Run with Docker Compose
docker-compose up

# Run tests
go test ./...

# Run integration tests
go test -tags=integration ./tests/integration/...
```

### Other Projects

Each project has its own README with specific development instructions:

- [event-bridge-kafka](./projects/event-bridge-kafka/README.md)
- [stream-anomaly-detector](./projects/stream-anomaly-detector/README.md)
- [stream-analytics-hub](./projects/stream-analytics-hub/README.md)
- [kafka-metrics-portal](./projects/kafka-metrics-portal/README.md)
- [log-replay-simulator](./projects/log-replay-simulator/README.md)

## Docker & Deployment

### Docker Compose

The main `docker-compose.yml` file orchestrates all services:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild specific service
docker-compose build distributed-tracing-system
docker-compose up -d distributed-tracing-system
```

### Individual Project Deployment

Each project can be deployed independently:

```bash
# Build Docker image
cd projects/distributed-tracing-system
docker build -t streamforge/distributed-tracing-system:latest .

# Run container
docker run -p 8082:8080 \
  -e JAEGER_ENDPOINT=http://jaeger:14268/api/traces \
  -e KAFKA_BROKERS=kafka:29092 \
  streamforge/distributed-tracing-system:latest
```

### Production Deployment

For production deployments, consider:

- **Kubernetes** - Container orchestration
- **Helm Charts** - Package management
- **CI/CD Pipelines** - Automated deployment
- **Service Mesh** - Istio for advanced traffic management
- **Secrets Management** - Vault or Kubernetes secrets

**Note:** Production deployment configurations are under active development. Some components may require additional configuration or may not be fully production-ready.

## Monitoring & Observability

### Prometheus Metrics

All services expose Prometheus metrics at `/metrics`:

```bash
# View metrics
curl http://localhost:8082/metrics
```

Key metrics:
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `traces_received_total` - Traces received
- `traces_processed_total` - Traces processed
- `kafka_messages_consumed_total` - Kafka messages consumed

### Grafana Dashboards

Access Grafana at http://localhost:3000

Pre-configured dashboards:
- **Service Overview** - High-level service metrics
- **Tracing Metrics** - Trace collection and processing
- **Kafka Metrics** - Message throughput and lag
- **System Health** - Infrastructure metrics

**Note:** Some dashboards may be in development or require manual configuration.

### Jaeger Tracing

Access Jaeger UI at http://localhost:16686

Features:
- Search traces by service, operation, or tags
- View trace timeline and spans
- Analyze service dependencies
- Performance analysis (latency percentiles)

### Logging

All services use structured JSON logging:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "service": "distributed-tracing-system",
  "message": "Trace processed successfully",
  "trace_id": "abc123",
  "span_id": "def456",
  "duration_ms": 45
}
```

View logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f distributed-tracing-system
```

## Security

### Best Practices

- **Non-root containers** - All services run as non-root users
- **Secrets management** - Environment variables for sensitive data
- **Network isolation** - Docker networks for service isolation
- **Health checks** - Container health monitoring
- **Resource limits** - CPU and memory constraints
- **Structured logging** - No sensitive data in logs

### Security Recommendations

For production:
- Use **HashiCorp Vault** or **AWS Secrets Manager** for secrets
- Enable **TLS/SSL** for all inter-service communication
- Implement **authentication and authorization** (JWT, OAuth2)
- Use **network policies** in Kubernetes
- Enable **audit logging** for compliance
- Regular **security scanning** of container images

**Note:** Security features are being continuously improved. Some security enhancements may be in development.

## Documentation

### Project Documentation

- [Development Guide](./docs/README.md) - Setup and development workflow
- [Architecture Documentation](./docs/README.md) - System design details
- [API Reference](./docs/README.md) - API documentation

### Project-Specific READMEs

Each project includes detailed documentation:

- [distributed-tracing-system](./projects/distributed-tracing-system/README.md)
- [event-bridge-kafka](./projects/event-bridge-kafka/README.md)
- [stream-anomaly-detector](./projects/stream-anomaly-detector/README.md)
- [stream-analytics-hub](./projects/stream-analytics-hub/README.md)
- [kafka-metrics-portal](./projects/kafka-metrics-portal/README.md)

**Note:** Documentation is continuously updated as the project evolves. Some sections may be incomplete or under review.

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific project
cd projects/distributed-tracing-system
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./tests/integration/...

# Run benchmarks
go test -bench=. ./...
```

### Test Strategy

- **Unit Tests** - Domain logic and utilities
- **Integration Tests** - Service integrations (Kafka, PostgreSQL, Jaeger)
- **E2E Tests** - Complete workflows
- **Performance Tests** - Load and stress testing

**Note:** Test coverage varies by project. Some projects may have incomplete test suites as they are under active development.

## Roadmap

### Phase 1: Foundation (Completed)
- [x] event-bridge-kafka
- [x] log-replay-simulator  
- [x] stream-analytics-hub
- [x] kafka-metrics-portal
- [x] stream-anomaly-detector

### Phase 2: Observability (In Progress)
- [x] distributed-tracing-system (core implementation)
- [ ] Enhanced trace storage and querying
- [ ] Trace correlation across services
- [ ] Performance optimization
- [ ] Additional exporters and integrations

### Phase 3: Intelligence (Planned)
- [ ] intelligent-alert-manager
- [ ] Advanced ML models
- [ ] Alert correlation and deduplication
- [ ] Predictive analytics

### Phase 4: User Experience (Planned)
- [ ] real-time-dashboard-builder
- [ ] Enhanced UI/UX
- [ ] Mobile-responsive dashboards
- [ ] Custom widget library

### Phase 5: Enterprise (Planned)
- [ ] multi-tenant-isolation-engine
- [ ] stream-data-validator
- [ ] Advanced security features
- [ ] Compliance and audit features
- [ ] Multi-region deployment support

**Note:** Roadmap items are subject to change based on development priorities and community feedback.

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** your changes (`git commit -m 'Add some AmazingFeature'`)
4. **Push** to the branch (`git push origin feature/AmazingFeature`)
5. **Open** a Pull Request

### Contribution Guidelines

- Follow the existing code style
- Write tests for new features
- Update documentation as needed
- Ensure all tests pass
- Follow semantic commit messages
- Check project status before contributing to ensure you're working on active components

### Report Issues

If you find a bug or have a suggestion, please [open an issue](https://github.com/your-org/stream-forge/issues).

**Note:** When reporting issues, please mention which project and version you're using, as development status varies across components.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

### Support Channels

- **Email**: support@streamforge.dev
- **Discussions**: [GitHub Discussions](https://github.com/your-org/stream-forge/discussions)
- **Issues**: [GitHub Issues](https://github.com/your-org/stream-forge/issues)
- **Documentation**: [Project Wiki](https://github.com/your-org/stream-forge/wiki)

### Use Cases

**Fintech**
- Real-time transaction monitoring
- ML-powered fraud detection
- Credit risk analysis

**IoT**
- Sensor data aggregation
- Device anomaly detection
- Energy efficiency optimization

**AI/ML**
- Production model monitoring
- Data drift detection
- Performance optimization

**Gaming**
- Player behavior analysis
- Bot and cheat detection
- Matchmaking optimization

---

<div align="center">

**Made with ❤️ by the StreamForge team**

[Star this project](https://github.com/your-org/stream-forge) • [Report bug](https://github.com/your-org/stream-forge/issues) • [Suggest feature](https://github.com/your-org/stream-forge/discussions)

</div>
