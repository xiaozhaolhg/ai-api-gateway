package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"github.com/ai-api-gateway/provider-service/internal/application"
	"github.com/ai-api-gateway/provider-service/internal/handler"
	"github.com/ai-api-gateway/provider-service/internal/infrastructure/adapter"
	"github.com/ai-api-gateway/provider-service/internal/infrastructure/repository"
	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	billingServiceAddr = "billing-service:50054"
	monitorServiceAddr = "monitor-service:50055"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Provider service listening on :50053")

	dbPath := "/data/provider.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Provider{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	providerRepo := repository.NewProviderRepository(db)
	adapterFactory := adapter.NewAdapterFactory()
	svc := application.NewService(providerRepo, adapterFactory, "default-crypto-key")

	h := handler.NewHandler(svc)

	log.Println("Registering callbacks with billing and monitor services...")
	time.Sleep(2 * time.Second)

	if err := registerCallback("billing-service", billingServiceAddr); err != nil {
		log.Printf("Warning: Failed to register with billing-service: %v", err)
	}
	if err := registerCallback("monitor-service", monitorServiceAddr); err != nil {
		log.Printf("Warning: Failed to register with monitor-service: %v", err)
	}
	log.Println("Callbacks registered (mock)")

	s := grpc.NewServer()
	providerv1.RegisterProviderServiceServer(s, h)

	log.Println("gRPC server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerCallback(serviceName, addr string) error {
	log.Printf("Registering callback with %s at %s...", serviceName, addr)
	log.Printf("Callback registered with %s (mock)", serviceName)
	return nil
}
