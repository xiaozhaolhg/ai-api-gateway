package testutils

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
)

// CreateBillingServer creates a test gRPC server for billing service
func CreateBillingServer() (billingv1.BillingServiceClient, net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, nil, nil, err
	}

	server := grpc.NewServer()
	mock := &mockBillingServer{
		usageRecords: make([]*billingv1.UsageRecord, 0),
	}
	billingv1.RegisterBillingServiceServer(server, mock)

	go func() {
		// Server.Serve will return error when listener is closed, which is expected
		server.Serve(listener)
	}()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		server.Stop()
		return nil, nil, nil, err
	}

	client := billingv1.NewBillingServiceClient(conn)
	return client, listener, server, nil
}

type mockBillingServer struct {
	billingv1.UnimplementedBillingServiceServer
	mu           sync.Mutex
	usageRecords []*billingv1.UsageRecord
}

func (m *mockBillingServer) RecordUsage(ctx context.Context, req *billingv1.RecordUsageRequest) (*commonv1.Empty, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	record := &billingv1.UsageRecord{
		UserId:           req.GetUserId(),
		GroupId:          req.GetGroupId(),
		ProviderId:       req.GetProviderId(),
		Model:            req.GetModel(),
		PromptTokens:     req.GetPromptTokens(),
		CompletionTokens: req.GetCompletionTokens(),
		Timestamp:        0, // Simplified for test
	}
	m.usageRecords = append(m.usageRecords, record)

	return &commonv1.Empty{}, nil
}

func (m *mockBillingServer) GetUsageRecords() []*billingv1.UsageRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.usageRecords
}

func (m *mockBillingServer) ClearRecords() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usageRecords = make([]*billingv1.UsageRecord, 0)
}
