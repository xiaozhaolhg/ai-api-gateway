package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Gateway service starting...")

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Gateway health check
	r.GET("/gateway/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Admin authentication endpoints (simple MVP implementation)
	admin := r.Group("/admin/auth")
	{
		admin.POST("/login", func(c *gin.Context) {
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "Invalid request"})
				return
			}
			// Simple MVP login - accept any credentials
			c.JSON(200, gin.H{
				"token": "mock-jwt-token-" + req.Username,
				"user": gin.H{
					"id":    "admin-user-id",
					"name":  req.Username,
					"email": "admin@example.com",
					"role":  "admin",
				},
			})
		})
		admin.POST("/logout", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Logged out successfully"})
		})
		admin.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"id":    "admin-user-id",
				"name":  "Admin",
				"email": "admin@example.com",
				"role":  "admin",
			})
		})
	}

	log.Printf("Gateway service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
