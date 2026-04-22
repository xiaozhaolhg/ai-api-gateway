package handler

import (
	"context"

	"github.com/ai-api-gateway/api/gen/billing/v1"
	"github.com/ai-api-gateway/billing-service/internal/application"
)

// Handler implements the BillingService gRPC interface
type Handler struct {
	billingv1.UnimplementedBillingServiceServer
	service *application.Service
}

// NewHandler creates a new billing service handler
func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// OnProviderResponse handles provider response callback
func (h *Handler) OnProviderResponse(ctx context.Context, req *billingv1.OnProviderResponseRequest) (*billingv1.Empty, error) {
	// Extract callback data
	callback := req.Callback

	// Record usage
	err := h.service.RecordUsage(callback.UserId, callback.ProviderId, callback.Model,
		callback.PromptTokens, callback.CompletionTokens)
	if err != nil {
		return nil, err
	}

	return &billingv1.Empty{}, nil
}

// RecordUsage records usage for a request
func (h *Handler) RecordUsage(ctx context.Context, req *billingv1.RecordUsageRequest) (*billingv1.Empty, error) {
	err := h.service.RecordUsage(req.UserId, req.ProviderId, req.Model,
		req.PromptTokens, req.CompletionTokens)
	if err != nil {
		return nil, err
	}

	return &billingv1.Empty{}, nil
}

// GetUsage retrieves usage records
func (h *Handler) GetUsage(ctx context.Context, req *billingv1.GetUsageRequest) (*billingv1.ListUsageResponse, error) {
	records, total, err := h.service.GetUsage(req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	// Convert domain entities to proto messages
	usageRecords := make([]*billingv1.UsageRecord, len(records))
	for i, record := range records {
		usageRecords[i] = &billingv1.UsageRecord{
			Id:               record.ID,
			UserId:           record.UserID,
			ProviderId:       record.ProviderID,
			Model:            record.Model,
			PromptTokens:     record.PromptTokens,
			CompletionTokens: record.CompletionTokens,
			Cost:             record.Cost,
			Timestamp:        record.Timestamp.Unix(),
		}
	}

	return &billingv1.ListUsageResponse{
		UsageRecords: usageRecords,
		Total:         int32(total),
	}, nil
}

// GetUsageAggregation retrieves aggregated usage statistics
func (h *Handler) GetUsageAggregation(ctx context.Context, req *billingv1.GetUsageAggregationRequest) (*billingv1.UsageAggregationResponse, error) {
	agg, err := h.service.GetUsageAggregation(req.UserId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &billingv1.UsageAggregationResponse{
		UserId:                agg.UserID,
		ProviderId:            agg.ProviderID,
		Model:                 agg.Model,
		TotalRequests:         agg.TotalRequests,
		TotalPromptTokens:     agg.TotalPromptTokens,
		TotalCompletionTokens: agg.TotalCompletionTokens,
		TotalCost:             agg.TotalCost,
		StartDate:             agg.StartDate,
		EndDate:               agg.EndDate,
	}, nil
}

// EstimateCost estimates the cost of a request
func (h *Handler) EstimateCost(ctx context.Context, req *billingv1.EstimateCostRequest) (*billingv1.CostEstimate, error) {
	cost, err := h.service.EstimateCost(req.ProviderId, req.Model,
		req.PromptTokens, req.CompletionTokens)
	if err != nil {
		return nil, err
	}

	return &billingv1.CostEstimate{
		EstimatedCost: cost,
	}, nil
}

// CheckBudget checks if a user has budget available
func (h *Handler) CheckBudget(ctx context.Context, req *billingv1.CheckBudgetRequest) (*billingv1.BudgetStatus, error) {
	status, err := h.service.CheckBudget(req.UserId)
	if err != nil {
		return nil, err
	}

	return &billingv1.BudgetStatus{
		CurrentSpend:    status.CurrentSpend,
		BudgetLimit:     status.BudgetLimit,
		Remaining:       status.Remaining,
		SoftCapExceeded: status.SoftCapExceeded,
		HardCapExceeded: status.HardCapExceeded,
	}, nil
}

// CreateBudget creates a new budget
func (h *Handler) CreateBudget(ctx context.Context, req *billingv1.CreateBudgetRequest) (*billingv1.Budget, error) {
	err := h.service.CreateBudget(req.Budget)
	if err != nil {
		return nil, err
	}

	return req.Budget, nil
}

// UpdateBudget updates an existing budget
func (h *Handler) UpdateBudget(ctx context.Context, req *billingv1.UpdateBudgetRequest) (*billingv1.Budget, error) {
	err := h.service.UpdateBudget(req.Budget)
	if err != nil {
		return nil, err
	}

	return req.Budget, nil
}

// DeleteBudget deletes a budget
func (h *Handler) DeleteBudget(ctx context.Context, req *billingv1.DeleteBudgetRequest) (*billingv1.Empty, error) {
	err := h.service.DeleteBudget(req.Id)
	if err != nil {
		return nil, err
	}

	return &billingv1.Empty{}, nil
}

// CreatePricingRule creates a new pricing rule
func (h *Handler) CreatePricingRule(ctx context.Context, req *billingv1.CreatePricingRuleRequest) (*billingv1.PricingRule, error) {
	err := h.service.CreatePricingRule(req.Rule)
	if err != nil {
		return nil, err
	}

	return req.Rule, nil
}

// UpdatePricingRule updates an existing pricing rule
func (h *Handler) UpdatePricingRule(ctx context.Context, req *billingv1.UpdatePricingRuleRequest) (*billingv1.PricingRule, error) {
	err := h.service.UpdatePricingRule(req.Rule)
	if err != nil {
		return nil, err
	}

	return req.Rule, nil
}

// DeletePricingRule deletes a pricing rule
func (h *Handler) DeletePricingRule(ctx context.Context, req *billingv1.DeletePricingRuleRequest) (*billingv1.Empty, error) {
	err := h.service.DeletePricingRule(req.Id)
	if err != nil {
		return nil, err
	}

	return &billingv1.Empty{}, nil
}

// GetBillingAccount retrieves a billing account
func (h *Handler) GetBillingAccount(ctx context.Context, req *billingv1.GetBillingAccountRequest) (*billingv1.BillingAccount, error) {
	account, err := h.service.GetBillingAccount(req.UserId)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// CreateBillingAccount creates a new billing account
func (h *Handler) CreateBillingAccount(ctx context.Context, req *billingv1.CreateBillingAccountRequest) (*billingv1.BillingAccount, error) {
	err := h.service.CreateBillingAccount(req.Account)
	if err != nil {
		return nil, err
	}

	return req.Account, nil
}

// UpdateBillingAccount updates an existing billing account
func (h *Handler) UpdateBillingAccount(ctx context.Context, req *billingv1.UpdateBillingAccountRequest) (*billingv1.BillingAccount, error) {
	err := h.service.UpdateBillingAccount(req.Account)
	if err != nil {
		return nil, err
	}

	return req.Account, nil
}

// GenerateInvoice generates an invoice
func (h *Handler) GenerateInvoice(ctx context.Context, req *billingv1.GenerateInvoiceRequest) (*billingv1.Invoice, error) {
	// Placeholder for MVP - invoice generation is Phase 3+
	return &billingv1.Invoice{}, nil
}

// GetInvoices retrieves invoices for a user
func (h *Handler) GetInvoices(ctx context.Context, req *billingv1.GetInvoicesRequest) (*billingv1.ListInvoicesResponse, error) {
	// Placeholder for MVP - invoice generation is Phase 3+
	return &billingv1.ListInvoicesResponse{}, nil
}
