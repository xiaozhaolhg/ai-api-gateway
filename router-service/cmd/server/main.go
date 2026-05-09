package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	"github.com/ai-api-gateway/router-service/internal/application"
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/handler"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/cache"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/config"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/repository"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	configPath := "configs/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&entity.RoutingRule{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repository
	ruleRepo := repository.NewRoutingRuleRepository(db)

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cache.Config{
		Address:    cfg.Cache.Redis.Address,
		Password:   cfg.Cache.Redis.Password,
		DB:         cfg.Cache.Redis.DB,
		DefaultTTL: cfg.Cache.Redis.TTLSeconds,
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without cache.", err)
		redisCache = nil
	}
	defer func() {
		if redisCache != nil {
			redisCache.Close()
		}
	}()

	// Initialize application service
	providerAddr := cfg.ProviderService.Host + ":" + cfg.ProviderService.Port
	if providerAddr == ":" {
		providerAddr = "localhost:50053"
	}
	appService, err := application.NewService(ruleRepo, redisCache, providerAddr)
	if err != nil {
		log.Fatalf("Failed to create application service: %v", err)
	}

	// Initialize gRPC handler
	grpcHandler := handler.NewHandler(appService)

	// Create gRPC server
	lis, err := net.Listen("tcp", cfg.Server.Host+":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	routerv1.RegisterRouterServiceServer(s, grpcHandler)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.GracefulStop()
		grpcHandler.Shutdown(ctx)
		log.Println("Server stopped")
	}()

	log.Printf("Router service listening on %s:%s", cfg.Server.Host, cfg.Server.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
