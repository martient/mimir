package observability

import (
	"fmt"
	"log"
	"time"
)

// Logger provides structured logging
type Logger struct {
	level string
}

// NewLogger creates a new logger
func NewLogger(cfg interface{}) (*Logger, error) {
	return &Logger{
		level: "info",
	}, nil
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] %s %v", msg, args)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, args)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] %s %v", msg, args)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, args)
}

// Metrics provides application metrics
type Metrics struct {
	requestsTotal   map[string]int
	requestDuration map[string][]time.Duration
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		requestsTotal:   make(map[string]int),
		requestDuration: make(map[string][]time.Duration),
	}
}

// RecordRequest records a request
func (m *Metrics) RecordRequest(path string, duration time.Duration) {
	m.requestsTotal[path]++
	m.requestDuration[path] = append(m.requestDuration[path], duration)
}

// PrometheusString returns metrics in Prometheus format
func (m *Metrics) PrometheusString() string {
	output := "# HELP mimir_requests_total Total number of requests\n"
	output += "# TYPE mimir_requests_total counter\n"

	for path, count := range m.requestsTotal {
		output += fmt.Sprintf("mimir_requests_total{path=\"%s\"} %d\n", path, count)
	}

	output += "\n# HELP mimir_request_duration_seconds Request duration in seconds\n"
	output += "# TYPE mimir_request_duration_seconds histogram\n"

	return output
}
