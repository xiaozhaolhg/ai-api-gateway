package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/util"
)

type JWTMiddleware struct {
	authClient *client.AuthClient
}

func NewJWTMiddleware(authClient *client.AuthClient) *JWTMiddleware {
	return &JWTMiddleware{authClient: authClient}
}

func (m *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(401, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		claims, err := util.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("userName", claims.Name)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func (m *JWTMiddleware) extractToken(c *gin.Context) string {
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		return token
	}

	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(403, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

func GetUserID(c *gin.Context) string {
	if id, exists := c.Get("userId"); exists {
		return id.(string)
	}
	return ""
}

func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("userRole"); exists {
		return role.(string)
	}
	return ""
}