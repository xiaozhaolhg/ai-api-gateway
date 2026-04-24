package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// RouteMiddleware resolves model to provider using router-service
type RouteMiddleware struct {
	routerClient *client.RouterClient
}

// NewRouteMiddleware creates a new route middleware
func NewRouteMiddleware(routerClient *client.RouterClient) *RouteMiddleware {
	return &RouteMiddleware{
		routerClient: routerClient,
	}
}

// Middleware returns the middleware function
func (m *RouteMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get model from request
		model := m.extractModel(r)
		if model == "" {
			http.Error(w, "Model not specified", http.StatusBadRequest)
			return
		}

		// Resolve route using router-service
		routeResult, err := m.routerClient.ResolveRoute(r.Context(), model)
		if err != nil {
			http.Error(w, "Failed to resolve route: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Add route result to context
		ctx := context.WithValue(r.Context(), "providerId", routeResult.ProviderId)
		ctx = context.WithValue(ctx, "adapterType", routeResult.AdapterType)
		ctx = context.WithValue(ctx, "fallbackProviderIds", routeResult.FallbackProviderIds)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractModel extracts the model name from the request
func (m *RouteMiddleware) extractModel(r *http.Request) string {
	// Try query parameter first
	model := r.URL.Query().Get("model")
	if model != "" {
		return model
	}

	// For POST requests, try to extract from request body
	// This is a simplified version - in production, you'd parse the JSON body
	if r.Method == http.MethodPost {
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err == nil {
			if model, ok := reqBody["model"].(string); ok {
				return model
			}
		}
	}

	return ""
}
