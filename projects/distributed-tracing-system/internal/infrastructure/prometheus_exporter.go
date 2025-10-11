package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// prometheusExporter implements the PrometheusExporter interface
type prometheusExporter struct {
	server           *http.Server
	registry         *prometheus.Registry
	tracesReceived   *prometheus.CounterVec
	tracesProcessed  *prometheus.CounterVec
	traceDuration    *prometheus.HistogramVec
	serviceLatency   *prometheus.HistogramVec
	errorRate        *prometheus.GaugeVec
}

// NewPrometheusExporter creates a new Prometheus exporter
func NewPrometheusExporter(port, path string) (domain.PrometheusExporter, error) {
	// Create custom registry
	registry := prometheus.NewRegistry()

	// Create metrics
	tracesReceived := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "traces_received_total",
			Help: "Total number of traces received",
		},
		[]string{"service", "operation"},
	)

	tracesProcessed := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "traces_processed_total",
			Help: "Total number of traces processed",
		},
		[]string{"service", "operation", "status"},
	)

	traceDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "trace_duration_seconds",
			Help:    "Duration of trace processing",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation"},
	)

	serviceLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_latency_seconds",
			Help:    "Service latency distribution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation"},
	)

	errorRate := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "error_rate",
			Help: "Error rate per service",
		},
		[]string{"service"},
	)

	// Register metrics
	registry.MustRegister(tracesReceived)
	registry.MustRegister(tracesProcessed)
	registry.MustRegister(traceDuration)
	registry.MustRegister(serviceLatency)
	registry.MustRegister(errorRate)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	exporter := &prometheusExporter{
		server:          server,
		registry:        registry,
		tracesReceived:  tracesReceived,
		tracesProcessed: tracesProcessed,
		traceDuration:   traceDuration,
		serviceLatency:  serviceLatency,
		errorRate:       errorRate,
	}

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Prometheus server error: %v\n", err)
		}
	}()

	return exporter, nil
}

// RecordTraceMetrics records metrics for a trace
func (pe *prometheusExporter) RecordTraceMetrics(trace *domain.Trace) error {
	if trace == nil {
		return fmt.Errorf("trace cannot be nil")
	}

	// Validate trace
	if err := pe.validateTrace(trace); err != nil {
		return fmt.Errorf("invalid trace: %w", err)
	}

	// Record metrics
	pe.tracesReceived.WithLabelValues(
		string(trace.Service),
		string(trace.Operation),
	).Inc()

	pe.tracesProcessed.WithLabelValues(
		string(trace.Service),
		string(trace.Operation),
		string(trace.Status),
	).Inc()

	pe.traceDuration.WithLabelValues(
		string(trace.Service),
		string(trace.Operation),
	).Observe(trace.Duration.Seconds())

	pe.serviceLatency.WithLabelValues(
		string(trace.Service),
		string(trace.Operation),
	).Observe(trace.Duration.Seconds())

	// Calculate error rate
	errorRate := 0.0
	if trace.Status == domain.TraceStatusError {
		errorRate = 1.0
	}
	pe.errorRate.WithLabelValues(string(trace.Service)).Set(errorRate)

	return nil
}

// validateTrace validates a trace before recording metrics
func (pe *prometheusExporter) validateTrace(trace *domain.Trace) error {
	if trace.Service == "" {
		return fmt.Errorf("service name is required")
	}
	if trace.Operation == "" {
		return fmt.Errorf("operation name is required")
	}
	if trace.Duration < 0 {
		return fmt.Errorf("duration cannot be negative")
	}
	return nil
}
