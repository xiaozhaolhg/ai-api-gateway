package handler

import (
	"context"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// BillingService defines the interface for billing operations
type BillingService interface {
	GetUsage(ctx context.Context, userID string, page, pageSize int32) (*UsageResp, error)
}

type UsageRecord struct {
	UserID           string  `json:"user_id"`
	Provider         string  `json:"provider"`
	Model            string  `json:"model"`
	PromptTokens     int64   `json:"prompt_tokens"`
	CompletionTokens int64   `json:"completion_tokens"`
	Cost             float64 `json:"cost"`
}

type UsageResp struct {
	Records []UsageRecord `json:"records"`
}

type AdminUsageHandler struct {
	billingClient *client.BillingClient
}

func NewAdminUsageHandler(billingClient *client.BillingClient) *AdminUsageHandler {
	return &AdminUsageHandler{billingClient: billingClient}
}

func (h *AdminUsageHandler) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*UsageResp, error) {
	resp, err := h.billingClient.GetUsage(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert client.UsageRecord to handler.UsageRecord
	records := make([]UsageRecord, len(resp.Records))
	for i, r := range resp.Records {
		records[i] = UsageRecord{
			UserID:           r.UserID,
			Provider:         r.Provider,
			Model:            r.Model,
			PromptTokens:     r.PromptTokens,
			CompletionTokens: r.CompletionTokens,
			Cost:             r.Cost,
		}
	}

	return &UsageResp{Records: records}, nil
}