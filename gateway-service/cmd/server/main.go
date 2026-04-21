package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Gateway service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
