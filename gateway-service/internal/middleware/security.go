package middleware

import (
	"net/http"
)

// SecurityMiddleware is a placeholder for security checks (MVP pass-through)
// In production, this would implement prompt injection detection, content filtering, etc.
type SecurityMiddleware struct{}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{}
}

// Middleware returns the middleware function (pass-through for MVP)
func (m *SecurityMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// MVP: Pass-through without security checks
		// In production, this would:
		// 1. Scan request body for prompt injection patterns
		// 2. Check against content moderation policies
		// 3. Validate request format and size limits
		// 4. Check for malicious patterns in headers
		// 5. Return 400 Bad Request if security checks fail

		next.ServeHTTP(w, r)
	})
}
