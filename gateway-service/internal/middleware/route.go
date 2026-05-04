package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

type RouteMiddleware struct {
	routerClient *client.RouterClient
}

func NewRouteMiddleware(routerClient *client.RouterClient) *RouteMiddleware {
	return &RouteMiddleware{
		routerClient: routerClient,
	}
}

func (m *RouteMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizedModels, _ := c.Get("authorizedModels")
		models := []string{}
		if authorizedModels != nil {
			models, _ = authorizedModels.([]string)
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		var req struct {
			Model string `json:"model"`
		}
		if err := json.Unmarshal(body, &req); err != nil || req.Model == "" {
			c.JSON(400, gin.H{"error": "Invalid request: model required"})
			c.Abort()
			return
		}

		result, err := m.routerClient.ResolveRoute(c.Request.Context(), req.Model, models)
		if err != nil {
			c.JSON(404, gin.H{"error": "Model not found: " + req.Model})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "providerId", result.ProviderID)
		ctx = context.WithValue(ctx, "fallbackProviderIds", result.FallbackProviderIDs)
		ctx = context.WithValue(ctx, "fallbackModels", result.FallbackModels)
		c.Request = c.Request.WithContext(ctx)
		c.Set("adapterType", result.AdapterType)
		c.Set("fallbackProviderIds", result.FallbackProviderIDs)
		c.Set("fallbackModels", result.FallbackModels)

		c.Next()
	}
}
