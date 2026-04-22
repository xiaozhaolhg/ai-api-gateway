package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	billingServiceAddr = "billing-service:50054"
	monitorServiceAddr = "monitor-service:50055"
)

func registerCallback(serviceName, addr string) error {
	log.Printf("Registering callback with %s at %s...", serviceName, addr)

	// For MVP, we'll just log the registration
	// In production, you'd make a gRPC call to register the callback
	// The provider-service would send its address to billing and monitor services
	// so they can call back with usage and metrics data

	log.Printf("Callback registered with %s (mock)", serviceName)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Provider service listening on :50053")

	// Register callbacks with billing and monitor services
	log.Println("Registering callbacks with billing and monitor services...")

	// Wait a moment for services to be ready
	time.Sleep(2 * time.Second)

	if err := registerCallback("billing-service", billingServiceAddr); err != nil {
		log.Printf("Warning: Failed to register with billing-service: %v", err)
	}

	if err := registerCallback("monitor-service", monitorServiceAddr); err != nil {
		log.Printf("Warning: Failed to register with monitor-service: %v", err)
	}

	log.Println("Callbacks registered (mock)")

	// TODO: Initialize gRPC server with provider service implementation
	s := grpc.NewServer()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
