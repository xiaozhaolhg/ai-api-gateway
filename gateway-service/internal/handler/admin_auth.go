package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/util"
)

type AdminAuthHandler struct {
	authServiceAddr string
}

func NewAdminAuthHandler(authServiceAddr string) *AdminAuthHandler {
	return &AdminAuthHandler{
		authServiceAddr: authServiceAddr,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles admin login
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// For MVP, we'll accept any username/password and return a token
	// In production, this would validate against the auth-service
	// For now, we'll create a mock user and generate a token
	
	// TODO: Implement proper password validation with auth-service
	// For now, we'll use a simple check
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a token for the user
	token, err := util.GenerateToken("admin-user-id", req.Username, "admin@example.com", "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    "admin-user-id",
			"name":  req.Username,
			"email": "admin@example.com",
			"role":  "admin",
		},
	})
}

// Logout handles admin logout
func (h *AdminAuthHandler) Logout(c *gin.Context) {
	// In a real implementation, we might invalidate the token
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetCurrentUser returns the current user info
func (h *AdminAuthHandler) GetCurrentUser(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	// Extract the token (remove "Bearer " prefix)
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Validate the token
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
