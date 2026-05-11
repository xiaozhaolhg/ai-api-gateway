package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	"github.com/ai-api-gateway/billing-service/internal/handler"
	"github.com/ai-api-gateway/billing-service/internal/application"
	"github.com/ai-api-gateway/billing-service/internal/infrastructure/repository"
	"github.com/ai-api-gateway/billing-service/internal/infrastructure/migration"
)

func main() {
	os.MkdirAll("/data", 0755)

	dbPath := "/data/billing.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := migration.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	usageRepo := repository.NewUsageRecordRepository(db)
	pricingRepo := repository.NewPricingRuleRepository(db)
	accountRepo := repository.NewBillingAccountRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)

	billingService := application.NewService(usageRepo, pricingRepo, accountRepo, budgetRepo)
	h := handler.NewHandler(billingService)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Billing service listening on :50054")

	s := grpc.NewServer()
	billingv1.RegisterBillingServiceServer(s, h)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
