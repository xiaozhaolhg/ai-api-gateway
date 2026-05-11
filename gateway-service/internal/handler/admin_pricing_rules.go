package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
)

// AdminPricingRulesHandler handles pricing rule management endpoints
type AdminPricingRulesHandler struct {
	billingClient *client.BillingClient
}

// NewAdminPricingRulesHandler creates a new admin pricing rules handler
func NewAdminPricingRulesHandler(billingClient *client.BillingClient) *AdminPricingRulesHandler {
	// Use billingv1 to avoid unused import error
	_ = billingv1.PricingRule{}
	
	return &AdminPricingRulesHandler{
		billingClient: billingClient,
	}
}

// PricingRuleResponse represents pricing rule response format for UI compatibility
type PricingRuleResponse struct {
	ID              string  `json:"id"`
	Model           string  `json:"model"`
	ProviderID      string  `json:"provider"`
	PromptPrice     float64 `json:"prompt_price"`
	CompletionPrice float64 `json:"completion_price"`
	Currency        string  `json:"currency"`
	EffectiveDate   string  `json:"effective_date"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ListPricingRules returns all pricing rules
func (h *AdminPricingRulesHandler) ListPricingRules(c *gin.Context) {
	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	resp, err := h.billingClient.ListPricingRules(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert proto pricing rules to UI-compatible format
	pricingRules := make([]PricingRuleResponse, len(resp.Rules))
	for i, rule := range resp.Rules {
		pricingRules[i] = PricingRuleResponse{
			ID:              rule.GetId(),
			Model:           rule.GetModel(),
			ProviderID:      rule.GetProviderId(),
			PromptPrice:     rule.GetPricePerPromptToken() * 1000, // Convert from per-token to per-1K
			CompletionPrice: rule.GetPricePerCompletionToken() * 1000, // Convert from per-token to per-1K
			Currency:        rule.GetCurrency(),
			EffectiveDate:   "", // TODO: Add effective_date to proto
			CreatedAt:       "", // TODO: Add created_at to proto
			UpdatedAt:       "", // TODO: Add updated_at to proto
		}
	}

	c.JSON(http.StatusOK, pricingRules)
}

// CreatePricingRule creates a new pricing rule
func (h *AdminPricingRulesHandler) CreatePricingRule(c *gin.Context) {
	var req struct {
		Model           string  `json:"model" binding:"required"`
		ProviderID      string  `json:"provider" binding:"required"`
		PromptPrice     float64 `json:"prompt_price" binding:"required"`
		CompletionPrice float64 `json:"completion_price" binding:"required"`
		Currency        string  `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	// Convert prices from per-1K to per-token for proto
	promptPricePerToken := req.PromptPrice / 1000
	completionPricePerToken := req.CompletionPrice / 1000

	resp, err := h.billingClient.CreatePricingRule(
		c.Request.Context(),
		req.Model,
		req.ProviderID,
		promptPricePerToken,
		completionPricePerToken,
		req.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pricingRule := PricingRuleResponse{
		ID:              resp.GetId(),
		Model:           resp.GetModel(),
		ProviderID:      resp.GetProviderId(),
		PromptPrice:     req.PromptPrice,
		CompletionPrice: req.CompletionPrice,
		Currency:        resp.GetCurrency(),
		EffectiveDate:   "", // TODO: Add effective_date to proto
		CreatedAt:       "", // TODO: Add created_at to proto
		UpdatedAt:       "", // TODO: Add updated_at to proto
	}

	c.JSON(http.StatusCreated, pricingRule)
}

// UpdatePricingRule updates an existing pricing rule
func (h *AdminPricingRulesHandler) UpdatePricingRule(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Model           string  `json:"model"`
		ProviderID      string  `json:"provider"`
		PromptPrice     float64 `json:"prompt_price"`
		CompletionPrice float64 `json:"completion_price"`
		Currency        string  `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	// Convert prices from per-1K to per-token for proto
	promptPricePerToken := req.PromptPrice / 1000
	completionPricePerToken := req.CompletionPrice / 1000

	resp, err := h.billingClient.UpdatePricingRule(
		c.Request.Context(),
		id,
		req.Model,
		req.ProviderID,
		promptPricePerToken,
		completionPricePerToken,
		req.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pricingRule := PricingRuleResponse{
		ID:              resp.GetId(),
		Model:           resp.GetModel(),
		ProviderID:      resp.GetProviderId(),
		PromptPrice:     req.PromptPrice,
		CompletionPrice: req.CompletionPrice,
		Currency:        resp.GetCurrency(),
		EffectiveDate:   "", // TODO: Add effective_date to proto
		CreatedAt:       "", // TODO: Add created_at to proto
		UpdatedAt:       "", // TODO: Add updated_at to proto
	}

	c.JSON(http.StatusOK, pricingRule)
}

// DeletePricingRule deletes a pricing rule
func (h *AdminPricingRulesHandler) DeletePricingRule(c *gin.Context) {
	id := c.Param("id")

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	if err := h.billingClient.DeletePricingRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pricing rule deleted"})
}
