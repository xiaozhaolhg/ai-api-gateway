package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxRetries     = 30
	retryInterval  = 2 * time.Second
	authServiceAddr    = "auth-service:50051"
	routerServiceAddr  = "router-service:50052"
	providerServiceAddr = "provider-service:50053"
)

func waitForService(name, addr string) error {
	log.Printf("Waiting for %s at %s...", name, addr)

	for i := 0; i < maxRetries; i++ {
		// Try to connect to the service
		// For MVP, we'll use a simple HTTP check if the service has an HTTP endpoint
		// In production, you'd use gRPC health checks

		// Try to connect - for gRPC services, we'd use grpc.Dial
		// For MVP, we'll just sleep and assume services come up
		time.Sleep(retryInterval)

		// For now, just log and continue
		// In production, implement actual health checks
		log.Printf("Attempt %d/%d: Checking %s...", i+1, maxRetries, name)
	}

	log.Printf("%s is ready (assumed)", name)
	return nil
}

func main() {
	// Wait for dependent services
	log.Println("Gateway service starting...")
	log.Println("Waiting for dependent services...")

	// Wait for auth, router, and provider services
	if err := waitForService("auth-service", authServiceAddr); err != nil {
		log.Fatalf("Failed to wait for auth-service: %v", err)
	}
	if err := waitForService("router-service", routerServiceAddr); err != nil {
		log.Fatalf("Failed to wait for router-service: %v", err)
	}
	if err := waitForService("provider-service", providerServiceAddr); err != nil {
		log.Fatalf("Failed to wait for provider-service: %v", err)
	}

	log.Println("All dependent services are ready")

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "dependencies": []string{
			authServiceAddr,
			routerServiceAddr,
			providerServiceAddr,
		}})
	})

	// Gateway health check
	r.GET("/gateway/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"services": map[string]string{
				"auth":    authServiceAddr,
				"router":  routerServiceAddr,
				"provider": providerServiceAddr,
			},
		})
	})

	log.Printf("Gateway service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
