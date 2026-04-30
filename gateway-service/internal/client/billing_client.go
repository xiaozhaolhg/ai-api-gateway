package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BillingClient wraps the billing-service gRPC client with lazy connection.
type BillingClient struct {
	address string
	client  billingv1.BillingServiceClient
	conn    *grpc.ClientConn
	mu      sync.RWMutex
}

// UsageResponse represents a usage response.
type UsageResponse struct {
	Records []UsageRecord `json:"records"`
}

// UsageRecord represents a single usage record.
type UsageRecord struct {
	UserID           string  `json:"user_id"`
	Provider         string  `json:"provider"`
	Model            string  `json:"model"`
	PromptTokens     int64   `json:"prompt_tokens"`
	CompletionTokens int64   `json:"completion_tokens"`
	Cost             float64 `json:"cost"`
}

// NewBillingClient creates a new billing service gRPC client with lazy connection.
func NewBillingClient(address string) (*BillingClient, error) {
	if address == "" {
		address = "localhost:50054"
	}
	return &BillingClient{
		address: address,
	}, nil
}

// getClient returns the gRPC client, initializing lazily if needed
func (c *BillingClient) getClient() (billingv1.BillingServiceClient, error) {
	c.mu.RLock()
	if c.client != nil {
		defer c.mu.RUnlock()
		return c.client, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.client != nil {
		return c.client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(GRPCInterceptor(DefaultRetryConfig())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to billing service: %w", err)
	}

	c.conn = conn
	c.client = billingv1.NewBillingServiceClient(conn)
	return c.client, nil
}

// Close closes the gRPC connection.
func (c *BillingClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUsage retrieves usage records from billing service.
func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*UsageResponse, error) {
	client, err := c.getClient()
	if err != nil {
		// Return empty response on connection error (graceful fallback)
		return &UsageResponse{Records: []UsageRecord{}}, nil
	}

	req := &billingv1.GetUsageRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.GetUsage(ctx, req)
	if err != nil {
		// Return empty response on error (graceful fallback)
		return &UsageResponse{Records: []UsageRecord{}}, nil
	}

	records := make([]UsageRecord, len(resp.Records))
	for i, r := range resp.Records {
		records[i] = UsageRecord{
			UserID:           r.UserId,
			Provider:         r.ProviderId,
			Model:            r.Model,
			PromptTokens:     r.PromptTokens,
			CompletionTokens: r.CompletionTokens,
			Cost:             r.Cost,
		}
	}

	return &UsageResponse{Records: records}, nil
}