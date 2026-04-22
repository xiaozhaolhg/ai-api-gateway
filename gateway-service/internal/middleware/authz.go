package middleware

import (
	"context"
	"net/http"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// AuthzMiddleware checks model authorization using auth-service
type AuthzMiddleware struct {
	authClient *client.AuthClient
}

// NewAuthzMiddleware creates a new authz middleware
func NewAuthzMiddleware(authClient *client.AuthClient) *AuthzMiddleware {
	return &AuthzMiddleware{
		authClient: authClient,
	}
}

// Middleware returns the middleware function
func (m *AuthzMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user info from context (set by AuthMiddleware)
		userID, ok := r.Context().Value("userId").(string)
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		groupIDs, _ := r.Context().Value("groupIds").([]string)

		// Get model from request (should be set by previous middleware or extracted from request body)
		model := r.URL.Query().Get("model")
		if model == "" {
			// Try to extract from request body for POST requests
			// For MVP, we'll skip body parsing and assume model is in query param or context
			http.Error(w, "Model not specified", http.StatusBadRequest)
			return
		}

		// Check model authorization with auth-service
		authResult, err := m.authClient.CheckModelAuthorization(r.Context(), userID, groupIDs, model)
		if err != nil {
			http.Error(w, "Failed to check authorization", http.StatusInternalServerError)
			return
		}

		if !authResult.Allowed {
			http.Error(w, "Model not authorized: "+authResult.Reason, http.StatusForbidden)
			return
		}

		// Add authorized models to context
		ctx := context.WithValue(r.Context(), "authorizedModels", authResult.AuthorizedModels)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
