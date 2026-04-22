package middleware

import (
	"log"
	"net/http"
	"time"
)

// LogMiddleware logs request metadata
type LogMiddleware struct{}

// NewLogMiddleware creates a new log middleware
func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

// Middleware returns the middleware function
func (m *LogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:      http.StatusOK,
		}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log request metadata
		duration := time.Since(start)
		m.logRequest(r, wrapped.statusCode, duration)
	})
}

// logRequest logs request metadata
func (m *LogMiddleware) logRequest(r *http.Request, statusCode int, duration time.Duration) {
	// Extract user ID from context if available
	userID, _ := r.Context().Value("userId").(string)
	if userID == "" {
		userID = "anonymous"
	}

	// Extract provider ID from context if available
	providerID, _ := r.Context().Value("providerId").(string)

	// Extract model from context if available
	model, _ := r.Context().Value("model").(string)

	// Log the request
	log.Printf(
		"method=%s path=%s status=%d duration=%s user_id=%s provider_id=%s model=%s",
		r.Method,
		r.URL.Path,
		statusCode,
		duration,
		userID,
		providerID,
		model,
	)
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
