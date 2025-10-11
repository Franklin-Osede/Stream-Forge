package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamforge/distributed-tracing-system/internal/config"
	"github.com/streamforge/distributed-tracing-system/internal/domain"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config       *config.Config
	traceService domain.TraceService
	router       *gin.Engine
	server       *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, traceService domain.TraceService) (*Server, error) {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
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

	s := &Server{
		config:       cfg,
		traceService: traceService,
		router:       router,
		server:       server,
	}

	// Setup routes
	s.setupRoutes()

	return s, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	log.Printf("Starting server on port %s", s.config.Server.Port)

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
	return s.server.Shutdown(shutdownCtx)
}

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
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
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "distributed-tracing-system",
	})
}

// searchTraces handles trace search requests
func (s *Server) searchTraces(c *gin.Context) {
	// Parse query parameters
	criteria := &domain.SearchCriteria{
		Limit:  10,
		Offset: 0,
	}

	// Get service filter
	if service := c.Query("service"); service != "" {
		serviceName := domain.ServiceName(service)
		criteria.Service = &serviceName
	}

	// Get operation filter
	if operation := c.Query("operation"); operation != "" {
		operationName := domain.OperationName(operation)
		criteria.Operation = &operationName
	}

	// Get limit
	if limit := c.Query("limit"); limit != "" {
		if l, err := parseInt(limit); err == nil && l > 0 {
			criteria.Limit = l
		}
	}

	// Get offset
	if offset := c.Query("offset"); offset != "" {
		if o, err := parseInt(offset); err == nil && o >= 0 {
			criteria.Offset = o
		}
	}

	// Search traces
	traces, err := s.traceService.SearchTraces(c.Request.Context(), criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"traces": traces,
		"total":  len(traces),
	})
}

// getTrace handles get trace by ID requests
func (s *Server) getTrace(c *gin.Context) {
	traceID := domain.TraceID(c.Param("id"))
	if traceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trace ID is required",
		})
		return
	}

	trace, err := s.traceService.GetTrace(c.Request.Context(), traceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if trace == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "trace not found",
		})
		return
	}

	c.JSON(http.StatusOK, trace)
}

// getServices handles get services requests
func (s *Server) getServices(c *gin.Context) {
	services, err := s.traceService.GetServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

// getOperations handles get operations requests
func (s *Server) getOperations(c *gin.Context) {
	service := domain.ServiceName(c.Param("service"))
	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "service name is required",
		})
		return
	}

	operations, err := s.traceService.GetOperations(c.Request.Context(), service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service":    service,
		"operations": operations,
	})
}

// getMetrics handles get metrics requests
func (s *Server) getMetrics(c *gin.Context) {
	metrics, err := s.traceService.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// Helper function to parse integers
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

