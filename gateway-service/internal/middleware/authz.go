package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/gin-gonic/gin"
)

type AuthzMiddleware struct {
	authClient *client.AuthClient
}

func NewAuthzMiddleware(authClient *client.AuthClient) *AuthzMiddleware {
	return &AuthzMiddleware{
		authClient: authClient,
	}
}

func (m *AuthzMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[DEBUG] Authz middleware called\n")
		userID, _ := c.Get("userId")
		if userID == nil {
			fmt.Printf("[DEBUG] User not authenticated\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}
		fmt.Printf("[DEBUG] User authenticated: %s\n", userID.(string))

		groupIDs, _ := c.Get("groupIds")
		groups := []string{}
		if groupIDs != nil {
			groups, _ = groupIDs.([]string)
		}

		model := c.Query("model")
		if model == "" && c.Request.Method == "POST" {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				c.Abort()
				return
			}
			// Reset body so downstream middlewares can read it
			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			var req struct {
				Model string `json:"model"`
			}
			if err := json.Unmarshal(body, &req); err == nil && req.Model != "" {
				model = req.Model
			}
		}

		if model == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Model not specified"})
			c.Abort()
			return
		}

		authResult, err := m.authClient.CheckModelAuthorization(c.Request.Context(), userID.(string), groups, model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check authorization"})
			c.Abort()
			return
		}

		if !authResult.Allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Model not authorized: " + authResult.Reason})
			c.Abort()
			return
		}

		c.Set("authorizedModels", authResult.AuthorizedModels)

		c.Next()
	}
}
