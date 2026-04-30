package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/pkg/cache"
	"github.com/gin-gonic/gin"
)

type RouteMiddleware struct {
	routerClient *client.RouterClient
	routeCache   *cache.Cache[string, string] // model -> providerID
}

func NewRouteMiddleware(routerClient *client.RouterClient) *RouteMiddleware {
	// Cache routing results for 10 minutes
	routeCache := cache.New[string, string](10 * time.Minute)
	routeCache.StartCleanup(2 * time.Minute)

	return &RouteMiddleware{
		routerClient: routerClient,
		routeCache:   routeCache,
	}
}

func (m *RouteMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := readBody(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		var req struct {
			Model string `json:"model"`
		}
		if err := json.Unmarshal(body, &req); err != nil || req.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: model required"})
			c.Abort()
			return
		}

		// Try cache first
		providerID, found := m.routeCache.Get(req.Model)
		if !found {
			// Cache miss - extract from model name
			providerID = extractProviderFromModel(req.Model)
			if providerID != "" {
				m.routeCache.Set(req.Model, providerID)
			}
		}

		if providerID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown model: " + req.Model})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "providerId", providerID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func extractProviderFromModel(model string) string {
	for i, c := range model {
		if c == ':' {
			return model[:i]
		}
	}
	return ""
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}
