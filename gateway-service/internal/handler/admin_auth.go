package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/util"
)

// AdminAuthHandler handles admin authentication via auth-service gRPC
type AdminAuthHandler struct {
	authClient *client.AuthClient
}

// NewAdminAuthHandler creates a new AdminAuthHandler
func NewAdminAuthHandler(authClient *client.AuthClient) *AdminAuthHandler {
	return &AdminAuthHandler{authClient: authClient}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles admin login via auth-service gRPC
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if h.authClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Auth service unavailable"})
		return
	}

	resp, err := h.authClient.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
		"user":  resp.User,
	})
}

// Logout handles admin logout
func (h *AdminAuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// CheckUsernameAvailability checks if a username is available
func (h *AdminAuthHandler) CheckUsernameAvailability(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if h.authClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Auth service unavailable"})
		return
	}

	available, err := h.authClient.CheckUsernameAvailability(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}

// GetCurrentUser returns the current user info from JWT claims
func (h *AdminAuthHandler) GetCurrentUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := util.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    claims.UserID,
		"name":  claims.Name,
		"email": claims.Email,
		"role":  claims.Role,
	})
}
