package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/taultek/mimir/internal/config"
	"github.com/taultek/mimir/internal/observability"
)

// Container holds shared dependencies
type Container struct {
	Config *config.Config
	Logger *observability.Logger
}

// NewContainer creates a new dependency container
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize logger
	logger, err := observability.NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	return &Container{
		Config: cfg,
		Logger: logger,
	}, nil
}

// Start starts the HTTP server
func (c *Container) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", c.healthHandler)

	// Metrics endpoint (simple version)
	mux.HandleFunc("/metrics", c.metricsHandler)

	// API routes (placeholder)
	mux.HandleFunc("/api/", c.apiHandler)

	// Simple logging middleware
	wrapped := loggingMiddleware(mux, c.Logger)
	mux.Handle("/", wrapped)

	// Start HTTP server
	addr := fmt.Sprintf(":%d", c.Config.Server.HTTPPort)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	c.Logger.Info("Starting Mimir Gateway", "port", c.Config.Server.HTTPPort)
	c.Logger.Info("Health check: http://localhost:%d/health", c.Config.Server.HTTPPort)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	c.Logger.Info("Shutting down Mimir Gateway...")
	return server.Shutdown(ctx)
}

// healthHandler returns health status
func (c *Container) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","version":"0.1.0"}`)
}

// metricsHandler returns metrics (placeholder)
func (c *Container) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# Mimir Metrics (placeholder)\n")
}

// apiHandler handles API requests (placeholder)
func (c *Container) apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Mimir API v0.1.0"}`)
}

// loggingMiddleware adds request logging
func loggingMiddleware(next http.Handler, logger *observability.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.status,
			"duration", duration.String(),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
