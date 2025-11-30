package server

import (
	"net/http"
	"time"

	"github.com/meysam81/oneoff/internal/metrics"
)

// MetricsMiddleware wraps HTTP handlers to collect request metrics
type MetricsMiddleware struct {
	collector metrics.Collector
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(collector metrics.Collector) *MetricsMiddleware {
	return &MetricsMiddleware{collector: collector}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware returns the middleware handler
func (m *MetricsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics for the /metrics endpoint itself
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status
		}

		// Serve the request
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start)
		path := normalizePath(r.URL.Path)
		m.collector.IncRequestsTotal(r.Method, path, wrapped.statusCode)
		m.collector.ObserveRequestDuration(r.Method, path, duration)
	})
}

// normalizePath reduces path cardinality by replacing IDs with placeholders
func normalizePath(path string) string {
	// Common API patterns with IDs
	patterns := map[string]string{
		"/api/jobs/":       "/api/jobs/:id",
		"/api/executions/": "/api/executions/:id",
		"/api/projects/":   "/api/projects/:id",
		"/api/tags/":       "/api/tags/:id",
		"/api/api-keys/":   "/api/api-keys/:id",
		"/api/webhooks/":   "/api/webhooks/:id",
	}

	for prefix, normalized := range patterns {
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			// Check for sub-routes like /api/jobs/:id/execute
			remaining := path[len(prefix):]
			for i, c := range remaining {
				if c == '/' {
					// Has sub-route, e.g., /api/jobs/123/execute -> /api/jobs/:id/execute
					return normalized + remaining[i:]
				}
			}
			return normalized
		}
	}

	return path
}
