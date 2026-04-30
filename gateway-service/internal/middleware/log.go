package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogMiddleware provides structured JSON logging with request IDs and sensitive data masking
type LogMiddleware struct {
	level         string
	maskSensitive bool
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   string            `json:"timestamp"`
	Level       string            `json:"level"`
	RequestID   string            `json:"request_id"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	StatusCode  int               `json:"status_code"`
	Duration    string            `json:"duration"`
	UserID      string            `json:"user_id,omitempty"`
	ProviderID  string            `json:"provider_id,omitempty"`
	Model       string            `json:"model,omitempty"`
	ClientIP    string            `json:"client_ip"`
	UserAgent   string            `json:"user_agent,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
}

// NewLogMiddleware creates a new structured JSON log middleware
func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{
		level:         "info",
		maskSensitive: true,
	}
}

// GinMiddleware returns the Gin-compatible middleware function
func (m *LogMiddleware) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestId", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)
		m.logRequest(c, requestID, duration)
	}
}

// logRequest creates and outputs a structured log entry
func (m *LogMiddleware) logRequest(c *gin.Context, requestID string, duration time.Duration) {
	// Extract values from context
	userID, _ := c.Get("userId")
	providerID, _ := c.Get("providerId")
	model, _ := c.Get("model")

	// Mask sensitive query parameters
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if m.isSensitiveParam(key) {
			queryParams[key] = "***"
		} else if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       m.determineLevel(c.Writer.Status()),
		RequestID:   requestID,
		Method:      c.Request.Method,
		Path:        c.Request.URL.Path,
		StatusCode:  c.Writer.Status(),
		Duration:    duration.String(),
		UserID:      m.toString(userID),
		ProviderID:  m.toString(providerID),
		Model:       m.toString(model),
		ClientIP:    c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		QueryParams: queryParams,
	}

	// Output as JSON
	jsonData, _ := json.Marshal(entry)
	gin.DefaultWriter.Write(jsonData)
	gin.DefaultWriter.Write([]byte("\n"))
}

// isSensitiveParam checks if a parameter contains sensitive data
func (m *LogMiddleware) isSensitiveParam(key string) bool {
	sensitiveParams := []string{
		"token", "password", "secret", "api_key", "apikey",
		"auth", "authorization", "cookie", "session",
	}
	lowerKey := strings.ToLower(key)
	for _, sensitive := range sensitiveParams {
		if strings.Contains(lowerKey, sensitive) {
			return true
		}
	}
	return false
}

// determineLevel maps HTTP status to log level
func (m *LogMiddleware) determineLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	default:
		return "info"
	}
}

// toString safely converts interface to string
func (m *LogMiddleware) toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// responseWriter wraps http.ResponseWriter to capture status code (kept for compatibility)
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
