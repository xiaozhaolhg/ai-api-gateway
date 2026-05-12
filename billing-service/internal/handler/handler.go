package handler

import (
	"context"
	"fmt"
	"time"

	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	"github.com/ai-api-gateway/billing-service/internal/application"
	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
)

type Handler struct {
	billingv1.UnimplementedBillingServiceServer
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) OnProviderResponse(ctx context.Context, req *commonv1.ProviderResponseCallback) (*commonv1.Empty, error) {
	if req == nil {
		return &commonv1.Empty{}, nil
	}
	if err := h.service.RecordUsage(req.GetUserId(), req.GetGroupId(), req.GetProviderId(), req.GetModel(),
		req.GetPromptTokens(), req.GetCompletionTokens()); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) RecordUsage(ctx context.Context, req *billingv1.RecordUsageRequest) (*commonv1.Empty, error) {
	if err := h.service.RecordUsage(req.GetUserId(), req.GetGroupId(), req.GetProviderId(), req.GetModel(),
		req.GetPromptTokens(), req.GetCompletionTokens()); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) GetUsage(ctx context.Context, req *billingv1.GetUsageRequest) (*billingv1.ListUsageResponse, error) {
	records, total, err := h.service.GetUsage(req.GetUserId(), int(req.GetPage()), int(req.GetPageSize()),
		req.GetStartTime(), req.GetEndTime())
	if err != nil {
		return nil, err
	}
	usageRecords := make([]*billingv1.UsageRecord, len(records))
	for i, record := range records {
		usageRecords[i] = &billingv1.UsageRecord{
			Id: record.ID, UserId: record.UserID, ProviderId: record.ProviderID,
			Model: record.Model, PromptTokens: record.PromptTokens,
			CompletionTokens: record.CompletionTokens, Cost: record.Cost,
			Timestamp: record.Timestamp.Unix(),
		}
	}
	return &billingv1.ListUsageResponse{Records: usageRecords, Total: int32(total)}, nil
}

func (h *Handler) GetUsageAggregation(ctx context.Context, req *billingv1.GetUsageAggregationRequest) (*billingv1.ListUsageAggregationResponse, error) {
	startDate := time.Unix(req.GetStartTime(), 0).Format("2006-01-02")
	endDate := time.Unix(req.GetEndTime(), 0).Format("2006-01-02")
	aggs, err := h.service.GetUsageAggregation(req.GetUserId(), startDate, endDate, req.GetGroupBy())
	if err != nil {
		return nil, err
	}
	aggregations := make([]*billingv1.UsageAggregation, len(aggs))
	for i, agg := range aggs {
		aggregations[i] = &billingv1.UsageAggregation{
			GroupKey:          agg.UserID + "|" + agg.ProviderID + "|" + agg.Model,
			TotalPromptTokens: agg.TotalPromptTokens, TotalCompletionTokens: agg.TotalCompletionTokens,
			TotalCost: agg.TotalCost, RequestCount: agg.TotalRequests,
		}
	}
	return &billingv1.ListUsageAggregationResponse{Aggregations: aggregations}, nil
}

func (h *Handler) EstimateCost(ctx context.Context, req *billingv1.EstimateCostRequest) (*billingv1.CostEstimate, error) {
	cost, err := h.service.EstimateCost(req.GetProviderId(), req.GetModel(), req.GetPromptTokens(), req.GetCompletionTokens())
	if err != nil {
		return nil, err
	}
	return &billingv1.CostEstimate{EstimatedCost: cost}, nil
}

func (h *Handler) CheckBudget(ctx context.Context, req *billingv1.CheckBudgetRequest) (*billingv1.BudgetStatus, error) {
	status, err := h.service.CheckBudget(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &billingv1.BudgetStatus{
		CurrentSpend: status.CurrentSpend, BudgetLimit: status.BudgetLimit,
		Remaining: status.Remaining, SoftCapExceeded: status.SoftCapExceeded,
		HardCapExceeded: status.HardCapExceeded,
	}, nil
}

func (h *Handler) CreateBudget(ctx context.Context, req *billingv1.CreateBudgetRequest) (*billingv1.Budget, error) {
	// Map proto to entity - convert percentage to absolute cap values
	subj := &entity.Budget{
		UserID:  req.GetUserId(),
		Limit:   req.GetLimit(),
		Period:  req.GetPeriod(),
		SoftCap: req.GetLimit() * req.GetSoftCapPct() / 100,
		HardCap: req.GetLimit() * req.GetHardCapPct() / 100,
		Status:  "active",
	}
	if err := h.service.CreateBudget(subj); err != nil {
		return nil, err
	}
	// Fetch created budget to get generated ID and timestamps
	budget, err := h.service.GetBudgetByUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return budgetToProto(budget), nil
}

func (h *Handler) UpdateBudget(ctx context.Context, req *billingv1.UpdateBudgetRequest) (*billingv1.Budget, error) {
	existing, err := h.service.GetBudgetByID(req.GetId())
	if err != nil {
		return nil, err
	}
	// Update fields
	existing.UserID = req.GetUserId()
	existing.Limit = req.GetLimit()
	existing.Period = req.GetPeriod()
	existing.SoftCap = req.GetLimit() * req.GetSoftCapPct() / 100
	existing.HardCap = req.GetLimit() * req.GetHardCapPct() / 100
	existing.Status = req.GetStatus()
	if err := h.service.UpdateBudget(existing); err != nil {
		return nil, err
	}
	updated, err := h.service.GetBudgetByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return budgetToProto(updated), nil
}

func (h *Handler) DeleteBudget(ctx context.Context, req *billingv1.DeleteBudgetRequest) (*commonv1.Empty, error) {
	if err := h.service.DeleteBudget(req.GetId()); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) ListBudgets(ctx context.Context, req *billingv1.ListBudgetsRequest) (*billingv1.ListBudgetsResponse, error) {
	budgets, total, err := h.service.ListBudgets(int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	protoBudgets := make([]*billingv1.Budget, len(budgets))
	for i, b := range budgets {
		protoBudgets[i] = budgetToProto(b)
	}
	return &billingv1.ListBudgetsResponse{Budgets: protoBudgets, Total: int32(total)}, nil
}

func (h *Handler) CreatePricingRule(ctx context.Context, req *billingv1.CreatePricingRuleRequest) (*billingv1.PricingRule, error) {
	// Map per-token prices to per-1K prices
	rule := &entity.PricingRule{
		Model:                req.GetModel(),
		ProviderID:           req.GetProviderId(),
		PromptPricePer1K:     req.GetPricePerPromptToken() * 1000,
		CompletionPricePer1K: req.GetPricePerCompletionToken() * 1000,
		Currency:             req.GetCurrency(),
	}
	if err := h.service.CreatePricingRule(rule); err != nil {
		return nil, err
	}
	// Fetch to get generated ID
	created, err := h.service.GetPricingRuleByID(rule.ID)
	if err != nil {
		return nil, err
	}
	return pricingRuleToProto(created), nil
}

func (h *Handler) UpdatePricingRule(ctx context.Context, req *billingv1.UpdatePricingRuleRequest) (*billingv1.PricingRule, error) {
	rule := &entity.PricingRule{
		ID:                   req.GetId(),
		Model:                req.GetModel(),
		ProviderID:           req.GetProviderId(),
		PromptPricePer1K:     req.GetPricePerPromptToken() * 1000,
		CompletionPricePer1K: req.GetPricePerCompletionToken() * 1000,
		Currency:             req.GetCurrency(),
	}
	if err := h.service.UpdatePricingRule(rule); err != nil {
		return nil, err
	}
	updated, err := h.service.GetPricingRuleByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return pricingRuleToProto(updated), nil
}

func (h *Handler) DeletePricingRule(ctx context.Context, req *billingv1.DeletePricingRuleRequest) (*commonv1.Empty, error) {
	if err := h.service.DeletePricingRule(req.GetId()); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) ListPricingRules(ctx context.Context, req *billingv1.ListPricingRulesRequest) (*billingv1.ListPricingRulesResponse, error) {
	rules, total, err := h.service.ListPricingRules(int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	protoRules := make([]*billingv1.PricingRule, len(rules))
	for i, r := range rules {
		protoRules[i] = pricingRuleToProto(r)
	}
	return &billingv1.ListPricingRulesResponse{Rules: protoRules, Total: int32(total)}, nil
}

func (h *Handler) GetBillingAccount(ctx context.Context, req *billingv1.GetBillingAccountRequest) (*billingv1.BillingAccount, error) {
	account, err := h.service.GetBillingAccountByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return billingAccountToProto(account), nil
}

func (h *Handler) GetBillingAccountByUser(ctx context.Context, req *billingv1.GetBillingAccountByUserRequest) (*billingv1.BillingAccount, error) {
	account, err := h.service.GetBillingAccount(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return billingAccountToProto(account), nil
}

func (h *Handler) CreateBillingAccount(ctx context.Context, req *billingv1.CreateBillingAccountRequest) (*billingv1.BillingAccount, error) {
	account := &entity.BillingAccount{
		UserID:         req.GetUserId(),
		Name:           req.GetName(),
		BillingContact: req.GetBillingContact(),
		Balance:        req.GetInitialCredit(),
		CreditBalance:  req.GetInitialCredit(),
		Currency:       "USD",
		Status:         "active",
	}
	if err := h.service.CreateBillingAccount(account); err != nil {
		return nil, err
	}
	created, err := h.service.GetBillingAccountByID(account.ID)
	if err != nil {
		return nil, err
	}
	return billingAccountToProto(created), nil
}

func (h *Handler) UpdateBillingAccount(ctx context.Context, req *billingv1.UpdateBillingAccountRequest) (*billingv1.BillingAccount, error) {
	existing, err := h.service.GetBillingAccountByID(req.GetId())
	if err != nil {
		return nil, err
	}
	if req.GetName() != "" {
		existing.Name = req.GetName()
	}
	if req.GetBillingContact() != "" {
		existing.BillingContact = req.GetBillingContact()
	}
	if req.GetStatus() != "" {
		existing.Status = req.GetStatus()
	}
	if req.GetBalanceUpdated() {
		existing.Balance = req.GetBalance()
	}
	if err := h.service.UpdateBillingAccount(existing); err != nil {
		return nil, err
	}
	updated, err := h.service.GetBillingAccountByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return billingAccountToProto(updated), nil
}

// Invoice handlers are Phase 3+
func (h *Handler) GenerateInvoice(ctx context.Context, req *billingv1.GenerateInvoiceRequest) (*billingv1.Invoice, error) {
	return nil, fmt.Errorf("not implemented")
}
func (h *Handler) GetInvoices(ctx context.Context, req *billingv1.GetInvoicesRequest) (*billingv1.ListInvoicesResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// Helper functions for entity-to-proto conversion

func budgetToProto(b *entity.Budget) *billingv1.Budget {
	if b == nil {
		return nil
	}
	var softPct, hardPct float64
	if b.Limit > 0 {
		softPct = b.SoftCap / b.Limit * 100
		hardPct = b.HardCap / b.Limit * 100
	}
	return &billingv1.Budget{
		Id:         b.ID,
		UserId:     b.UserID,
		Limit:      b.Limit,
		Period:     b.Period,
		SoftCapPct: softPct,
		HardCapPct: hardPct,
		Status:     b.Status,
	}
}

func pricingRuleToProto(r *entity.PricingRule) *billingv1.PricingRule {
	if r == nil {
		return nil
	}
	return &billingv1.PricingRule{
		Id:                      r.ID,
		Model:                   r.Model,
		ProviderId:              r.ProviderID,
		PricePerPromptToken:     r.PromptPricePer1K / 1000,
		PricePerCompletionToken: r.CompletionPricePer1K / 1000,
		Currency:                r.Currency,
	}
}

func billingAccountToProto(a *entity.BillingAccount) *billingv1.BillingAccount {
	if a == nil {
		return nil
	}
	return &billingv1.BillingAccount{
		Id:             a.ID,
		UserId:         a.UserID,
		Name:           a.Name,
		BillingContact: a.BillingContact,
		Balance:        a.Balance,
		CreditBalance:  a.CreditBalance,
		Currency:       a.Currency,
		Status:         a.Status,
	}
}
