package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header
		apiKey := m.extractAPIKey(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		// Validate API key with auth-service
		userIdentity, err := m.authClient.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

	c.Set("userId", userIdentity.UserId)
	c.Set("role", userIdentity.Role)
	c.Set("groupIds", userIdentity.GroupIds)
	c.Set("scopes", userIdentity.Scopes)

	c.Next()
	}
}

func (m *AuthMiddleware) extractAPIKey(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Support "Bearer <token>" format
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	// Support raw API key
	return authHeader
}
