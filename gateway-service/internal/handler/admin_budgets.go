package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
)

// AdminBudgetsHandler handles budget management endpoints
type AdminBudgetsHandler struct {
	billingClient *client.BillingClient
}

// NewAdminBudgetsHandler creates a new admin budgets handler
func NewAdminBudgetsHandler(billingClient *client.BillingClient) *AdminBudgetsHandler {
	// Use billingv1 to avoid unused import error
	_ = billingv1.Budget{}
	
	return &AdminBudgetsHandler{
		billingClient: billingClient,
	}
}

// BudgetResponse represents budget response format for UI compatibility
type BudgetResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Scope         string  `json:"scope"`
	ScopeID       string  `json:"scope_id"`
	Limit         float64 `json:"limit"`
	CurrentSpend  float64 `json:"current_spend"`
	Period        string  `json:"period"`
	SoftCapPct    float64 `json:"soft_cap_pct"`
	HardCapPct    float64 `json:"hard_cap_pct"`
	Status        string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ListBudgets returns all budgets
func (h *AdminBudgetsHandler) ListBudgets(c *gin.Context) {
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

	resp, err := h.billingClient.ListBudgets(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert proto budgets to UI-compatible format
	budgets := make([]BudgetResponse, len(resp.Budgets))
	for i, budget := range resp.Budgets {
		budgets[i] = BudgetResponse{
			ID:            budget.GetId(),
			Name:          budget.GetUserId() + " Budget", // Generate name from user_id
			Scope:         "user",
			ScopeID:       budget.GetUserId(),
			Limit:         budget.GetLimit(),
			CurrentSpend:  0, // TODO: Calculate from usage records
			Period:        budget.GetPeriod(),
			SoftCapPct:    budget.GetSoftCapPct(),
			HardCapPct:    budget.GetHardCapPct(),
			Status:        budget.GetStatus(),
			CreatedAt:      "", // TODO: Add created_at to proto
			UpdatedAt:      "", // TODO: Add updated_at to proto
		}
	}

	c.JSON(http.StatusOK, budgets)
}

// CreateBudget creates a new budget
func (h *AdminBudgetsHandler) CreateBudget(c *gin.Context) {
	var req struct {
		UserID      string  `json:"user_id" binding:"required"`
		Limit       float64 `json:"limit" binding:"required"`
		Period      string  `json:"period" binding:"required"`
		SoftCapPct  float64 `json:"soft_cap_pct"`
		HardCapPct  float64 `json:"hard_cap_pct"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	resp, err := h.billingClient.CreateBudget(
		c.Request.Context(),
		req.UserID,
		req.Limit,
		req.Period,
		req.SoftCapPct,
		req.HardCapPct,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budget := BudgetResponse{
		ID:            resp.GetId(),
		Name:          req.UserID + " Budget",
		Scope:         "user",
		ScopeID:       req.UserID,
		Limit:         resp.GetLimit(),
		CurrentSpend:  0,
		Period:        resp.GetPeriod(),
		SoftCapPct:    resp.GetSoftCapPct(),
		HardCapPct:    resp.GetHardCapPct(),
		Status:        resp.GetStatus(),
		CreatedAt:      "", // TODO: Add created_at to proto
		UpdatedAt:      "", // TODO: Add updated_at to proto
	}

	c.JSON(http.StatusCreated, budget)
}

// UpdateBudget updates an existing budget
func (h *AdminBudgetsHandler) UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		UserID      string  `json:"user_id"`
		Limit       float64 `json:"limit"`
		Period      string  `json:"period"`
		SoftCapPct  float64 `json:"soft_cap_pct"`
		HardCapPct  float64 `json:"hard_cap_pct"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	resp, err := h.billingClient.UpdateBudget(
		c.Request.Context(),
		id,
		req.UserID,
		req.Limit,
		req.Period,
		req.SoftCapPct,
		req.HardCapPct,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budget := BudgetResponse{
		ID:            resp.GetId(),
		Name:          req.UserID + " Budget",
		Scope:         "user",
		ScopeID:       req.UserID,
		Limit:         resp.GetLimit(),
		CurrentSpend:  0,
		Period:        resp.GetPeriod(),
		SoftCapPct:    resp.GetSoftCapPct(),
		HardCapPct:    resp.GetHardCapPct(),
		Status:        resp.GetStatus(),
		CreatedAt:      "", // TODO: Add created_at to proto
		UpdatedAt:      "", // TODO: Add updated_at to proto
	}

	c.JSON(http.StatusOK, budget)
}

// DeleteBudget deletes a budget
func (h *AdminBudgetsHandler) DeleteBudget(c *gin.Context) {
	id := c.Param("id")

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	if err := h.billingClient.DeleteBudget(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "budget deleted"})
}
