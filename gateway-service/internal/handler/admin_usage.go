package handler

import (
	"context"
	"time"

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
	Timestamp        string  `json:"timestamp"`
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

// GetUsage retrieves usage records for a user with optional date range filtering.
// If startTime or endTime is 0, that bound is not applied.
func (h *AdminUsageHandler) GetUsage(ctx context.Context, userID string, page, pageSize int32, startTime, endTime int64) (*UsageResp, error) {
	resp, err := h.billingClient.GetUsage(ctx, userID, page, pageSize, startTime, endTime)
	if err != nil {
		return nil, err
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
			Timestamp:        time.Unix(r.Timestamp, 0).Format(time.RFC3339),
		}
	}

	return &UsageResp{Records: records}, nil
}