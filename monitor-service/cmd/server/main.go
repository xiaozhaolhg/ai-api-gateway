package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	monitorv1 "github.com/ai-api-gateway/api/gen/monitor/v1"
	"github.com/ai-api-gateway/monitor-service/internal/handler"
	"github.com/ai-api-gateway/monitor-service/internal/application"
	"github.com/ai-api-gateway/monitor-service/internal/infrastructure/repository"
	"github.com/ai-api-gateway/monitor-service/internal/infrastructure/migration"
)

func main() {
	os.MkdirAll("/data", 0755)

	dbPath := "/data/monitor.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := migration.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	metricRepo := repository.NewMetricRepository(db)
	alertRuleRepo := repository.NewAlertRuleRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	healthRepo := repository.NewProviderHealthRepository(db)

	monitorService := application.NewService(metricRepo, alertRuleRepo, alertRepo, healthRepo)
	h := handler.NewHandler(monitorService)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Monitor service listening on :50055")

	s := grpc.NewServer()
	monitorv1.RegisterMonitorServiceServer(s, h)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}