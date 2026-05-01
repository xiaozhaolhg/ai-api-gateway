package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/handler"
	"github.com/ai-api-gateway/gateway-service/internal/infrastructure/config"
	"github.com/ai-api-gateway/gateway-service/internal/middleware"
	"github.com/ai-api-gateway/gateway-service/internal/util"
)

var (
	authClient     *client.AuthClient
	routerClient   *client.RouterClient
	providerClient *client.ProviderClient
	billingClient  *client.BillingClient
)

func main() {
	log.Println("Gateway service starting...")

	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Printf("Warning: failed to load config: %v", err)
		cfg = &config.Config{
			Server: config.ServerConfig{Port: "8080", Host: "0.0.0.0"},
		}
	}

	// Initialize clients with lazy connection
	var initErr error
	authClient, initErr = client.NewAuthClient(cfg.AuthService.Address)
	if initErr != nil {
		log.Printf("Warning: auth client initialization failed: %v", initErr)
	}

	routerClient, initErr = client.NewRouterClient(cfg.RouterService.Address)
	if initErr != nil {
		log.Printf("Warning: router client initialization failed: %v", initErr)
	}

	providerClient, initErr = client.NewProviderClient(cfg.ProviderService.Address)
	if initErr != nil {
		log.Printf("Warning: provider client initialization failed: %v", initErr)
	}

	billingClient, initErr = client.NewBillingClient(cfg.BillingService.Address)
	if initErr != nil {
		log.Printf("Warning: billing client initialization failed: %v", initErr)
	}

	authMiddleware := middleware.NewAuthMiddleware(authClient)
	authzMiddleware := middleware.NewAuthzMiddleware(authClient)
	routeMiddleware := middleware.NewRouteMiddleware(routerClient)
	proxyMiddleware := middleware.NewProxyMiddleware(providerClient, billingClient)

	r := gin.Default()

	// Add logging middleware
	logMiddleware := middleware.NewLogMiddleware()
	r.Use(logMiddleware.GinMiddleware())

	// CORS middleware
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

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(authClient, routerClient, providerClient, billingClient)
	modelsHandler := handler.NewModelsHandler(providerClient)
	adminUsageHandler := handler.NewAdminUsageHandler(billingClient)

	// Simple liveness check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Deep health check with dependencies
	r.GET("/gateway/health", healthHandler.GatewayHealth)

	// Models aggregation endpoint
	r.GET("/gateway/models", modelsHandler.ListModels)

	r.POST("/admin/auth/login", handleLogin)
	r.POST("/admin/auth/register", handleRegister)
	r.POST("/admin/auth/logout", handleLogout)

	admin := r.Group("/admin/auth")
	admin.Use(jwtAuthMiddleware())
	{
		admin.GET("/me", handleGetCurrentUser)

		// User management
		admin.GET("/users", handleListUsers)
		admin.POST("/users", handleCreateUser)
		admin.PUT("/users/:id", handleUpdateUser)
		admin.DELETE("/users/:id", handleDeleteUser)

		// API key management
		admin.GET("/api-keys/:user_id", handleListAPIKeys)
		admin.POST("/api-keys", handleCreateAPIKey)
		admin.DELETE("/api-keys/:id", handleDeleteAPIKey)

		// Group management
		admin.GET("/groups", handleListGroups)
		admin.POST("/groups", handleCreateGroup)
		admin.PUT("/groups/:id", handleUpdateGroup)
		admin.DELETE("/groups/:id", handleDeleteGroup)
		admin.POST("/groups/:id/members", handleAddUserToGroup)
		admin.DELETE("/groups/:id/members/:user_id", handleRemoveUserFromGroup)

		// Permission management
		admin.GET("/permissions", handleListPermissions)
		admin.POST("/permissions", handleGrantPermission)
		admin.DELETE("/permissions/:id", handleRevokePermission)

		// Usage (billing) - using adminUsageHandler from main branch
		admin.GET("/usage", func(c *gin.Context) {
			handleGetUsage(c, adminUsageHandler)
		})

	}

	providerHandler := handler.NewAdminProvidersHandler(cfg.ProviderService.Address, cfg.RouterService.Address)
	r.POST("/admin/providers", providerHandler.CreateProvider)
	r.GET("/admin/providers", providerHandler.ListProviders)
	r.PUT("/admin/providers/:id", providerHandler.UpdateProvider)
	r.DELETE("/admin/providers/:id", providerHandler.DeleteProvider)
	r.GET("/admin/providers/:id/health", providerHandler.HealthCheck)

	// Add error handling middleware (must be after logging to capture status codes)
	r.Use(middleware.NewErrorMiddleware().Middleware())

	v1 := r.Group("/v1")
	{
		chat := v1.Group("/chat/completions")
		chat.Use(
			authMiddleware.Middleware(),
			authzMiddleware.Middleware(),
			routeMiddleware.Middleware(),
			wrapHTTPMiddleware(proxyMiddleware.Middleware),
		)
		chat.POST("", func(c *gin.Context) {})

		models := v1.Group("/models")
		models.Use(authMiddleware.Middleware())
		models.GET("", modelsHandler.ListModels)

		v1.GET("/providers", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"providers": []string{"ollama", "opencode_zen"}})
		})
	}

	v1Auth := v1.Group("/auth")
	v1Auth.Use(jwtAuthMiddleware())
	{
		v1Auth.POST("/api-keys", handleCreateUserAPIKey)
		v1Auth.GET("/api-keys", handleListUserAPIKeys)
		v1Auth.DELETE("/api-keys/:id", handleDeleteUserAPIKey)
	}

	// Setup HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Gateway service listening on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close gRPC connections
	if authClient != nil {
		authClient.Close()
	}
	if routerClient != nil {
		routerClient.Close()
	}
	if providerClient != nil {
		providerClient.Close()
	}
	if billingClient != nil {
		billingClient.Close()
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
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

func wrapHTTPMiddleware(m func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		restOfChain := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})
		wrapped := m(restOfChain)
		wrapped.ServeHTTP(c.Writer, c.Request)

		if c.Writer.Written() {
			c.Abort()
		}
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

// --- Auth handlers ---

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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("[DEBUG] Login error for %s: %v", req.Email, err)
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	setAuthCookie(c, resp.Token)
	c.JSON(200, gin.H{
		"token": resp.Token,
		"user":  resp.User,
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Register(ctx, req.Username, req.Email, req.Name, req.Password, req.Role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, resp.Token)
	c.JSON(200, gin.H{
		"token": resp.Token,
		"user":  resp.User,
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
		"name":  c.GetString("userName"),
		"email": c.GetString("userEmail"),
		"role":  c.GetString("userRole"),
	})
}

// --- User management handlers (wired to auth-service gRPC) ---

func handleListUsers(c *gin.Context) {
	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	resp, err := authClient.ListUsers(context.Background(), page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"users": resp.Users, "total": resp.Total})
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

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	user, err := authClient.CreateUser(context.Background(), req.Name, req.Email, req.Role, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
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

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	user, err := authClient.UpdateUser(context.Background(), id, req.Name, req.Email, req.Role, req.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func handleDeleteUser(c *gin.Context) {
	id := c.Param("id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.DeleteUser(context.Background(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user deleted"})
}

// --- API Key management handlers ---

func handleListAPIKeys(c *gin.Context) {
	userID := c.Param("user_id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	resp, err := authClient.ListAPIKeys(context.Background(), userID, page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"api_keys": resp.ApiKeys, "total": resp.Total})
}

func handleCreateAPIKey(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	resp, err := authClient.CreateAPIKey(context.Background(), req.UserID, req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"api_key_id": resp.ApiKeyId, "api_key": resp.ApiKey})
}

func handleDeleteAPIKey(c *gin.Context) {
	id := c.Param("id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.DeleteAPIKey(context.Background(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "api key deleted"})
}

// --- Group management handlers ---

func handleListGroups(c *gin.Context) {
	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	resp, err := authClient.ListGroups(context.Background(), page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"groups": resp.Groups, "total": resp.Total})
}

func handleCreateGroup(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		ParentGroupID string `json:"parent_group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	group, err := authClient.CreateGroup(context.Background(), req.Name, req.ParentGroupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, group)
}

func handleUpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name          string `json:"name"`
		ParentGroupID string `json:"parent_group_id"`
	}
	c.ShouldBindJSON(&req)

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	group, err := authClient.UpdateGroup(context.Background(), id, req.Name, req.ParentGroupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, group)
}

func handleDeleteGroup(c *gin.Context) {
	id := c.Param("id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.DeleteGroup(context.Background(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "group deleted"})
}

func handleAddUserToGroup(c *gin.Context) {
	groupID := c.Param("id")
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.AddUserToGroup(context.Background(), req.UserID, groupID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user added to group"})
}

func handleRemoveUserFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.Param("user_id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.RemoveUserFromGroup(context.Background(), userID, groupID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user removed from group"})
}

// --- Permission management handlers ---

func handleListPermissions(c *gin.Context) {
	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	groupID := c.Query("group_id")
	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	resp, err := authClient.ListPermissions(context.Background(), groupID, page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"permissions": resp.Permissions, "total": resp.Total})
}

func handleGrantPermission(c *gin.Context) {
	var req struct {
		GroupID      string `json:"group_id" binding:"required"`
		ResourceType string `json:"resource_type" binding:"required"`
		ResourceID   string `json:"resource_id" binding:"required"`
		Action       string `json:"action" binding:"required"`
		Effect       string `json:"effect"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	permission, err := authClient.GrantPermission(context.Background(), req.GroupID, req.ResourceType, req.ResourceID, req.Action)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, permission)
}

func handleRevokePermission(c *gin.Context) {
	id := c.Param("id")

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	if err := authClient.RevokePermission(context.Background(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "permission revoked"})
}

func handleCreateUserAPIKey(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(401, gin.H{"error": "authorization required"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	resp, err := authClient.CreateAPIKey(context.Background(), userID, req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"api_key_id": resp.ApiKeyId,
		"api_key":     resp.ApiKey,
		"name":         req.Name,
	})
}

func handleListUserAPIKeys(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(401, gin.H{"error": "authorization required"})
		return
	}

	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	resp, err := authClient.ListAPIKeys(context.Background(), userID, page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	type apiKey struct {
		ID        string `json:"api_key_id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
	}
	keys := make([]apiKey, len(resp.ApiKeys))
	for i, k := range resp.ApiKeys {
		keys[i] = apiKey{
			ID:        k.Id,
			Name:      k.Name,
			CreatedAt: time.Unix(k.CreatedAt, 0).Format(time.RFC3339),
		}
	}

	c.JSON(200, gin.H{"api_keys": keys, "total": resp.Total})
}

func handleDeleteUserAPIKey(c *gin.Context) {
	userID := c.GetString("userId")
	keyID := c.Param("id")

	if userID == "" {
		c.JSON(401, gin.H{"error": "authorization required"})
		return
	}

	if authClient == nil {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}

	err := authClient.DeleteAPIKey(context.Background(), keyID)
	if err != nil {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(200, gin.H{"message": "api key deleted"})
}

// --- Usage handler (using adminUsageHandler from main branch) ---

func handleGetUsage(c *gin.Context, h *handler.AdminUsageHandler) {
	userID := c.GetString("userId")
	if userID == "" {
		userID = "anonymous"
	}

	page := int32(1)
	pageSize := int32(20)

	resp, err := h.GetUsage(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		log.Printf("Error getting usage: %v", err)
		// Return empty response on error (graceful fallback)
		c.JSON(200, gin.H{"usage": []gin.H{}, "error": err.Error()})
		return
	}

	usage := make([]gin.H, len(resp.Records))
	for i, r := range resp.Records {
		usage[i] = gin.H{
			"user_id":           r.UserID,
			"provider":          r.Provider,
			"model":             r.Model,
			"prompt_tokens":     r.PromptTokens,
			"completion_tokens": r.CompletionTokens,
			"cost":              r.Cost,
		}
	}

	c.JSON(200, gin.H{"usage": usage})
}

// --- Providers (still mock for now) ---

func handleListProviders(c *gin.Context) {
	c.JSON(200, gin.H{"providers": []gin.H{
		{"id": "ollama", "name": "Ollama", "enabled": true},
		{"id": "opencode_zen", "name": "OpenCode Zen", "enabled": true},
	}})
}
