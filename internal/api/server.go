package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"psql-mcp-registry/internal/instance_manager"

	"github.com/gin-gonic/gin"
)

// APIServer represents the HTTP API server
type APIServer struct {
	manager instance_manager.Manager
	router  *gin.Engine
	server  *http.Server
	port    string
}

// NewAPIServer creates a new HTTP API server instance
func NewAPIServer(manager instance_manager.Manager, port string) *APIServer {
	// Set Gin mode from environment (defaults to release mode)
	if gin.Mode() == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	apiServer := &APIServer{
		manager: manager,
		router:  router,
		port:    port,
	}

	// Register routes
	apiServer.registerRoutes()

	return apiServer
}

// registerRoutes sets up all API routes
func (s *APIServer) registerRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.HealthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		v1.POST("/instances", s.RegisterInstance)
	}
}

// Run starts the HTTP server
func (s *APIServer) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("HTTP server shutdown error: %w", err)
		}
		return nil
	case err := <-errChan:
		return err
	}
}

