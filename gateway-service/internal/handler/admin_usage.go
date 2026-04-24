package handler

import "context"

type BillingService interface {
	GetUsage(ctx context.Context, userID string, page, pageSize int) (*UsageResp, error)
}

type UsageRecord struct {
	UserID    string `json:"user_id"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	Cost    float64 `json:"cost"`
}

type UsageResp struct {
	Records []UsageRecord `json:"records"`
}

type AdminUsageHandler struct {
	svc BillingService
}

func NewAdminUsageHandler(svc BillingService) *AdminUsageHandler {
	return &AdminUsageHandler{svc: svc}
}

func (h *AdminUsageHandler) GetUsage(userID string, page, pageSize int) (*UsageResp, error) {
	return h.svc.GetUsage(context.Background(), userID, page, pageSize)
}