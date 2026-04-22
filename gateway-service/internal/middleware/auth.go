package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// AuthMiddleware validates API keys using auth-service
type AuthMiddleware struct {
	authClient *client.AuthClient
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authClient *client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Middleware returns the middleware function
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract API key from Authorization header
		apiKey := m.extractAPIKey(r)
		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		// Validate API key with auth-service
		userIdentity, err := m.authClient.ValidateAPIKey(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		// Add user identity to context
		ctx := context.WithValue(r.Context(), "userId", userIdentity.UserId)
		ctx = context.WithValue(ctx, "role", userIdentity.Role)
		ctx = context.WithValue(ctx, "groupIds", userIdentity.GroupIds)
		ctx = context.WithValue(ctx, "scopes", userIdentity.Scopes)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractAPIKey extracts the API key from the Authorization header
func (m *AuthMiddleware) extractAPIKey(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Support "Bearer <token>" format
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Support raw API key
	return authHeader
}
