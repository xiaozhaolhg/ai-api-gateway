package client

import (
	"context"
	"fmt"
)

type BillingClient struct{}

func NewBillingClient(address string) (*BillingClient, error) {
	return &BillingClient{}, nil
}

func (c *BillingClient) Close() error {
	return nil
}

func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*UsageResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

type UsageResponse struct {
	Records []UsageRecord `json:"records"`
}

type UsageRecord struct {
	UserID    string `json:"user_id"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	Cost      float64 `json:"cost"`
}