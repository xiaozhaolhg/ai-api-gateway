package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/handler"
	"github.com/ai-api-gateway/gateway-service/internal/util"
)

var authClient *client.AuthClient

func main() {
	log.Println("Gateway service starting...")

	var err error
	authAddress := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authAddress == "" {
		authAddress = "localhost:50051"
	}
	authClient, err = client.NewAuthClient(authAddress)
	if err != nil {
		log.Printf("Warning: failed to connect to auth service: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/gateway/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	r.POST("/admin/auth/login", handleLogin)
	r.POST("/admin/auth/register", handleRegister)
	r.POST("/admin/auth/logout", handleLogout)

	admin := r.Group("/admin/auth")
	admin.Use(jwtAuthMiddleware())
	{
		admin.GET("/me", handleGetCurrentUser)
		admin.GET("/users", handleListUsers)
		admin.POST("/users", handleCreateUser)
		admin.PUT("/users/:id", handleUpdateUser)
		admin.DELETE("/users/:id", handleDeleteUser)
		admin.GET("/usage", handleGetUsage)
	}

	providerHandler := handler.NewAdminProvidersHandler()
	r.POST("/admin/providers", providerHandler.CreateProvider)
	r.GET("/admin/providers", providerHandler.ListProviders)
	r.PUT("/admin/providers/:id", providerHandler.UpdateProvider)
	r.DELETE("/admin/providers/:id", providerHandler.DeleteProvider)
	r.GET("/admin/providers/:id/health", providerHandler.HealthCheck)

	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Chat completions"})
		})
		v1.GET("/models", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"models": []string{"ollama:llama2", "opencode_zen:gpt-4"}})
		})
		v1.GET("/providers", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"providers": []string{"ollama", "opencode_zen"}})
		})
	}

	log.Printf("Gateway service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
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

func extractToken(c *gin.Context) string {
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		return token
	}
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return authHeader
}

func handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	resp, err := authClient.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	setAuthCookie(c, resp.Token)
	c.JSON(200, gin.H{
		"token": resp.Token,
		"user": resp.User,
	})
}

func handleRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if req.Username == "" && req.Email == "" {
		c.JSON(400, gin.H{"error": "username or email required"})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	resp, err := authClient.Register(context.Background(), req.Username, req.Email, req.Name, req.Password, req.Role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, resp.Token)
	c.JSON(200, gin.H{
		"token": resp.Token,
		"user": resp.User,
	})
}

func handleLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		Path:   "/admin",
		MaxAge: -1,
	})
	c.JSON(200, gin.H{"message": "logged out"})
}

func handleGetCurrentUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"id":    c.GetString("userId"),
		"name": c.GetString("userName"),
		"email": c.GetString("userEmail"),
		"role": c.GetString("userRole"),
	})
}

func setAuthCookie(c *gin.Context, token string) {
	maxAge := 24 * 60 * 60
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/admin",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
	})
}

func handleListUsers(c *gin.Context) {
	c.JSON(200, gin.H{"users": []gin.H{
		{"id": "user-1", "name": "Admin", "email": "admin@example.com", "role": "admin"},
	}})
}

func handleCreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(201, gin.H{"id": "user-new", "name": req.Name, "email": req.Email, "role": req.Role})
}

func handleUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}
	c.ShouldBindJSON(&req)
	c.JSON(200, gin.H{"id": id, "name": req.Name, "email": req.Email, "role": req.Role})
}

func handleDeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "user deleted"})
}

func handleListProviders(c *gin.Context) {
	c.JSON(200, gin.H{"providers": []gin.H{
		{"id": "ollama", "name": "Ollama", "enabled": true},
		{"id": "opencode_zen", "name": "OpenCode Zen", "enabled": true},
	}})
}

func handleGetUsage(c *gin.Context) {
	c.JSON(200, gin.H{"usage": []gin.H{
		{"period": "24h", "tokens": 15000, "cost": 0.15},
	}})
}