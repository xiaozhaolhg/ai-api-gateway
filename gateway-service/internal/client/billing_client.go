package client

import (
	"context"
	"fmt"

	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BillingClient wraps the billing-service gRPC client
type BillingClient struct {
	client billingv1.BillingServiceClient
	conn   *grpc.ClientConn
}

// NewBillingClient creates a new billing service client
func NewBillingClient(address string) (*BillingClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to billing service: %w", err)
	}

	return &BillingClient{
		client: billingv1.NewBillingServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *BillingClient) Close() error {
	return c.conn.Close()
}

// GetUsage retrieves usage records
func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*billingv1.ListUsageResponse, error) {
	req := &billingv1.GetUsageRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.GetUsage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage: %w", err)
	}

	return resp, nil
}
