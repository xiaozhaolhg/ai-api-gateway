package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/auth-service/internal/handler"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/repository"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/migration"
)

func main() {
	os.MkdirAll("/data", 0755)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth service listening on :50051")

	dbPath := "/data/auth.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if _, err := os.Create(dbPath); err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migration.Migrate("/data/auth.db"); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	authService := application.NewAuthService(userRepo, apiKeyRepo)
	h := handler.NewHandler(authService, userRepo, apiKeyRepo)

	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, h)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}