package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/config"
	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/streamforge/distributed-tracing-system/internal/telemetry"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ServerWithTelemetry represents the HTTP server with OpenTelemetry instrumentation
type ServerWithTelemetry struct {
	config           *config.Config
	traceService     domain.TraceService
	telemetryManager *telemetry.TelemetryManager
	router           *gin.Engine
	server           *http.Server
}

// NewServerWithTelemetry creates a new server instance with telemetry
func NewServerWithTelemetry(cfg *config.Config, traceService domain.TraceService, telemetryManager *telemetry.TelemetryManager) (*ServerWithTelemetry, error) {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add telemetry middleware
	httpMiddleware := telemetry.NewHTTPMiddleware(
		telemetryManager.GetTracer(),
		telemetryManager.GetMeter(),
	)
	router.Use(httpMiddleware.GinMiddleware())

	// Add standard Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	s := &ServerWithTelemetry{
		config:           cfg,
		traceService:     traceService,
		telemetryManager: telemetryManager,
		router:           router,
		server:           server,
	}

	// Setup routes
	s.setupRoutes()

	return s, nil
}

// Start starts the server
func (s *ServerWithTelemetry) Start(ctx context.Context) error {
	log.Printf("Starting server with telemetry on port %s", s.config.Server.Port)

	// Start server in goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	
	// Shutdown telemetry
	if err := s.telemetryManager.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown telemetry: %v", err)
	}

	return s.server.Shutdown(shutdownCtx)
}

// setupRoutes sets up the HTTP routes with telemetry
func (s *ServerWithTelemetry) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Trace routes
		traces := v1.Group("/traces")
		{
			traces.GET("/search", s.searchTraces)
			traces.GET("/:id", s.getTrace)
		}

		// Service routes
		services := v1.Group("/services")
		{
			services.GET("", s.getServices)
			services.GET("/:service/operations", s.getOperations)
		}

		// Metrics routes
		metrics := v1.Group("/metrics")
		{
			metrics.GET("", s.getMetrics)
		}
	}
}

// healthCheck handles health check requests
func (s *ServerWithTelemetry) healthCheck(c *gin.Context) {
	// Create a span for health check
	_, span := s.telemetryManager.StartSpan(c.Request.Context(), "health-check")
	defer span.End()

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "distributed-tracing-system",
		"version":   "1.0.0",
	})
}

// searchTraces handles trace search requests
func (s *ServerWithTelemetry) searchTraces(c *gin.Context) {
	// Create a span for search operation
	ctx, span := s.telemetryManager.StartSpan(c.Request.Context(), "search-traces")
	defer span.End()

	// Parse query parameters
	criteria := &domain.SearchCriteria{
		Limit:  10,
		Offset: 0,
	}

	// Get service filter
	if service := c.Query("service"); service != "" {
		serviceName := domain.ServiceName(service)
		criteria.Service = &serviceName
		span.SetAttributes(attribute.String("search.service", service))
	}

	// Get operation filter
	if operation := c.Query("operation"); operation != "" {
		operationName := domain.OperationName(operation)
		criteria.Operation = &operationName
		span.SetAttributes(attribute.String("search.operation", operation))
	}

	// Get limit
	if limit := c.Query("limit"); limit != "" {
		if l, err := parseIntTelemetry(limit); err == nil && l > 0 {
			criteria.Limit = l
			span.SetAttributes(attribute.Int("search.limit", l))
		}
	}

	// Get offset
	if offset := c.Query("offset"); offset != "" {
		if o, err := parseIntTelemetry(offset); err == nil && o >= 0 {
			criteria.Offset = o
			span.SetAttributes(attribute.Int("search.offset", o))
		}
	}

	// Search traces
	traces, err := s.traceService.SearchTraces(ctx, criteria)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	span.SetAttributes(attribute.Int("search.results_count", len(traces)))
	span.SetStatus(codes.Ok, "Search completed successfully")

	c.JSON(http.StatusOK, gin.H{
		"traces": traces,
		"total":  len(traces),
	})
}

// getTrace handles get trace by ID requests
func (s *ServerWithTelemetry) getTrace(c *gin.Context) {
	// Create a span for get trace operation
	ctx, span := s.telemetryManager.StartSpan(c.Request.Context(), "get-trace")
	defer span.End()

	traceID := domain.TraceID(c.Param("id"))
	if traceID == "" {
		span.SetStatus(codes.Error, "Trace ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trace ID is required",
		})
		return
	}

	span.SetAttributes(attribute.String("trace.id", string(traceID)))

	trace, err := s.traceService.GetTrace(ctx, traceID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if trace == nil {
		span.SetStatus(codes.Error, "Trace not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "trace not found",
		})
		return
	}

	span.SetAttributes(
		attribute.String("trace.service", string(trace.Service)),
		attribute.String("trace.operation", string(trace.Operation)),
		attribute.String("trace.status", string(trace.Status)),
	)
	span.SetStatus(codes.Ok, "Trace retrieved successfully")

	c.JSON(http.StatusOK, trace)
}

// getServices handles get services requests
func (s *ServerWithTelemetry) getServices(c *gin.Context) {
	// Create a span for get services operation
	ctx, span := s.telemetryManager.StartSpan(c.Request.Context(), "get-services")
	defer span.End()

	services, err := s.traceService.GetServices(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	span.SetAttributes(attribute.Int("services.count", len(services)))
	span.SetStatus(codes.Ok, "Services retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

// getOperations handles get operations requests
func (s *ServerWithTelemetry) getOperations(c *gin.Context) {
	// Create a span for get operations operation
	ctx, span := s.telemetryManager.StartSpan(c.Request.Context(), "get-operations")
	defer span.End()

	service := domain.ServiceName(c.Param("service"))
	if service == "" {
		span.SetStatus(codes.Error, "Service name is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "service name is required",
		})
		return
	}

	span.SetAttributes(attribute.String("service.name", string(service)))

	operations, err := s.traceService.GetOperations(ctx, service)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	span.SetAttributes(attribute.Int("operations.count", len(operations)))
	span.SetStatus(codes.Ok, "Operations retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"service":    service,
		"operations": operations,
	})
}

// getMetrics handles get metrics requests
func (s *ServerWithTelemetry) getMetrics(c *gin.Context) {
	// Create a span for get metrics operation
	ctx, span := s.telemetryManager.StartSpan(c.Request.Context(), "get-metrics")
	defer span.End()

	metrics, err := s.traceService.GetMetrics(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	span.SetStatus(codes.Ok, "Metrics retrieved successfully")

	c.JSON(http.StatusOK, metrics)
}

// Helper function to parse integers
func parseIntTelemetry(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
