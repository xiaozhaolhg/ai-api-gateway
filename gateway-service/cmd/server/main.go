package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
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
		admin.GET("/providers", handleListProviders)
		admin.GET("/usage", handleGetUsage)
	}

	// Add routes that match UI expectations
	adminUI := r.Group("/admin")
	adminUI.Use(jwtAuthMiddleware())
	{
		adminUI.GET("/providers", handleListProviders)
		adminUI.POST("/providers", handleCreateProvider)
		adminUI.PUT("/providers/:id", handleUpdateProvider)
		adminUI.DELETE("/providers/:id", handleDeleteProvider)
		adminUI.GET("/users", handleListUsers)
		adminUI.POST("/users", handleCreateUser)
		adminUI.PUT("/users/:id", handleUpdateUser)
		adminUI.DELETE("/users/:id", handleDeleteUser)
		adminUI.GET("/usage", handleGetUsage)
		adminUI.GET("/api-keys/:userId", handleListAPIKeys)
		adminUI.POST("/api-keys", handleCreateAPIKey)
		adminUI.DELETE("/api-keys/:id", handleDeleteAPIKey)
		adminUI.GET("/routing-rules", handleListRoutingRules)
		adminUI.POST("/routing-rules", handleCreateRoutingRule)
		adminUI.PUT("/routing-rules/:id", handleUpdateRoutingRule)
		adminUI.DELETE("/routing-rules/:id", handleDeleteRoutingRule)
		adminUI.GET("/groups", handleListGroups)
		adminUI.POST("/groups", handleCreateGroup)
		adminUI.PUT("/groups/:id", handleUpdateGroup)
		adminUI.DELETE("/groups/:id", handleDeleteGroup)
		adminUI.GET("/permissions", handleListPermissions)
		adminUI.POST("/permissions", handleCreatePermission)
		adminUI.PUT("/permissions/:id", handleUpdatePermission)
		adminUI.DELETE("/permissions/:id", handleDeletePermission)
		adminUI.GET("/budgets", handleListBudgets)
		adminUI.POST("/budgets", handleCreateBudget)
		adminUI.PUT("/budgets/:id", handleUpdateBudget)
		adminUI.DELETE("/budgets/:id", handleDeleteBudget)
		adminUI.GET("/pricing-rules", handleListPricingRules)
		adminUI.POST("/pricing-rules", handleCreatePricingRule)
		adminUI.PUT("/pricing-rules/:id", handleUpdatePricingRule)
		adminUI.DELETE("/pricing-rules/:id", handleDeletePricingRule)
		adminUI.GET("/alert-rules", handleListAlertRules)
		adminUI.POST("/alert-rules", handleCreateAlertRule)
		adminUI.PUT("/alert-rules/:id", handleUpdateAlertRule)
		adminUI.DELETE("/alert-rules/:id", handleDeleteAlertRule)
		adminUI.GET("/alerts", handleListAlerts)
		adminUI.PUT("/alerts/:id/acknowledge", handleAcknowledgeAlert)
		adminUI.GET("/health", handleGetProviderHealth)
	}

	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Chat completions"})
		})
		v1.GET("/models", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"models": []string{"ollama:llama2", "opencode_zen:gpt-4", "local-ollama:qwen3.5:9b", "local-ollama:gemma4:e4b"}})
		})
		v1.GET("/providers", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"providers": []string{"ollama", "opencode_zen", "local-ollama"}})
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
	c.JSON(200, []gin.H{
		{"id": "user-1", "name": "Admin", "email": "admin@example.com", "role": "admin", "status": "active", "created_at": "2026-04-27T00:00:00Z"},
		{"id": "ac5451299fe5dea53b96bf53d72ff5dd", "name": "admin", "email": "admin@abc.com", "role": "admin", "status": "active", "created_at": "2026-04-27T13:14:18Z"},
		{"id": "2046007d99dc380466be33ec9f80083b", "name": "Super Admin", "email": "superadmin@example.com", "role": "admin", "status": "active", "created_at": "2026-04-27T13:29:12Z"},
	})
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
	c.JSON(200, []gin.H{
		{"id": "ollama", "name": "Ollama", "type": "ollama", "base_url": "http://localhost:11434", "enabled": true, "status": "active", "models": []string{"llama2"}, "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
		{"id": "opencode_zen", "name": "OpenCode Zen", "type": "openai", "base_url": "https://opencode.ai/zen", "enabled": true, "status": "active", "models": []string{"gpt-4"}, "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
		{"id": "local-ollama", "name": "Local Ollama", "type": "ollama", "base_url": "http://host.docker.internal:11434", "enabled": true, "status": "active", "models": []string{"qwen3.5:9b", "gemma4:e4b"}, "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateProvider(c *gin.Context) {
	var req struct {
		Name     string   `json:"name"`
		Type     string   `json:"type"`
		BaseURL  string   `json:"base_url"`
		Models   []string `json:"models"`
		Status   string   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(201, gin.H{
		"id": "provider-" + c.GetString("userId"),
		"name": req.Name,
		"type": req.Type,
		"base_url": req.BaseURL,
		"models": req.Models,
		"status": req.Status,
		"enabled": true,
		"created_at": "2026-04-27T00:00:00Z",
		"updated_at": "2026-04-27T00:00:00Z",
	})
}

func handleUpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name     string   `json:"name"`
		Type     string   `json:"type"`
		BaseURL  string   `json:"base_url"`
		Models   []string `json:"models"`
		Status   string   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(200, gin.H{
		"id": id,
		"name": req.Name,
		"type": req.Type,
		"base_url": req.BaseURL,
		"models": req.Models,
		"status": req.Status,
		"enabled": true,
		"updated_at": "2026-04-27T00:00:00Z",
	})
}

func handleDeleteProvider(c *gin.Context) {
	c.JSON(200, gin.H{"message": "provider deleted"})
}

func handleGetUsage(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "usage-1", "user_id": "user-1", "model": "local-ollama:qwen3.5:9b", "provider": "local-ollama", "prompt_tokens": 10000, "completion_tokens": 5000, "total_tokens": 15000, "cost": 0.15, "timestamp": "2026-04-27T12:00:00Z"},
	})
}

// Additional handlers for admin UI endpoints
func handleListAPIKeys(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "key-1", "user_id": "user-1", "name": "Default Key", "scopes": []string{"read", "write"}, "created_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateAPIKey(c *gin.Context) {
	c.JSON(200, gin.H{"api_key_id": "key-new", "api_key": "sk-new-key-12345"})
}

func handleDeleteAPIKey(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListRoutingRules(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "rule-1", "model_pattern": "gpt-*", "provider": "openai", "priority": 1, "fallback_chain": []string{}, "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateRoutingRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": "rule-new", "model_pattern": "qwen-*", "provider": "local-ollama", "priority": 1, "fallback_chain": []string{}, "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdateRoutingRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "model_pattern": "updated", "provider": "local-ollama", "priority": 1, "fallback_chain": []string{}, "status": "active", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeleteRoutingRule(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListGroups(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "group-1", "name": "Admin Group", "description": "Administrators", "member_count": 2, "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateGroup(c *gin.Context) {
	c.JSON(200, gin.H{"id": "group-new", "name": "New Group", "description": "New group", "member_count": 0, "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdateGroup(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "name": "Updated Group", "description": "Updated description", "member_count": 1, "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeleteGroup(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListPermissions(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "perm-1", "group_id": "group-1", "model_pattern": "*", "effect": "allow", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreatePermission(c *gin.Context) {
	c.JSON(200, gin.H{"id": "perm-new", "group_id": "group-1", "model_pattern": "gpt-*", "effect": "allow", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdatePermission(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "group_id": "group-1", "model_pattern": "updated", "effect": "allow", "status": "active", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeletePermission(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListBudgets(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "budget-1", "name": "Default Budget", "scope": "global", "scope_id": "", "limit": 100.0, "current_spend": 15.5, "period": "monthly", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateBudget(c *gin.Context) {
	c.JSON(200, gin.H{"id": "budget-new", "name": "New Budget", "scope": "global", "scope_id": "", "limit": 50.0, "current_spend": 0.0, "period": "monthly", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdateBudget(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "name": "Updated Budget", "scope": "global", "scope_id": "", "limit": 75.0, "current_spend": 20.0, "period": "monthly", "status": "active", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeleteBudget(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListPricingRules(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "price-1", "model": "gpt-4", "provider": "openai", "prompt_price": 0.03, "completion_price": 0.06, "currency": "USD", "effective_date": "2026-04-27T00:00:00Z", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreatePricingRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": "price-new", "model": "qwen3.5:9b", "provider": "local-ollama", "prompt_price": 0.01, "completion_price": 0.02, "currency": "USD", "effective_date": "2026-04-27T00:00:00Z", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdatePricingRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "model": "updated", "provider": "local-ollama", "prompt_price": 0.015, "completion_price": 0.025, "currency": "USD", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeletePricingRule(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListAlertRules(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "alert-1", "name": "Budget Alert", "metric": "spend", "condition": ">", "threshold": 80.0, "channel": "email", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"},
	})
}

func handleCreateAlertRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": "alert-new", "name": "New Alert", "metric": "tokens", "condition": ">", "threshold": 1000.0, "channel": "webhook", "status": "active", "created_at": "2026-04-27T00:00:00Z", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleUpdateAlertRule(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "name": "Updated Alert", "metric": "tokens", "condition": ">", "threshold": 1500.0, "channel": "email", "status": "active", "updated_at": "2026-04-27T00:00:00Z"})
}

func handleDeleteAlertRule(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleListAlerts(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "alert-1", "rule_id": "alert-1", "severity": "warning", "status": "active", "triggered_at": "2026-04-27T12:00:00Z", "description": "Budget at 80%", "acknowledged_at": ""},
	})
}

func handleAcknowledgeAlert(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func handleGetProviderHealth(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "local-ollama", "name": "Local Ollama", "status": "healthy", "latency_ms": 50, "error_rate": 0.0, "last_check": "2026-04-27T15:00:00Z"},
		{"id": "ollama", "name": "Ollama", "status": "healthy", "latency_ms": 100, "error_rate": 0.0, "last_check": "2026-04-27T15:00:00Z"},
		{"id": "opencode_zen", "name": "OpenCode Zen", "status": "healthy", "latency_ms": 200, "error_rate": 0.0, "last_check": "2026-04-27T15:00:00Z"},
	})
}