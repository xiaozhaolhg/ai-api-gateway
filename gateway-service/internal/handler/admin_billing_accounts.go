package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

type AdminBillingAccountsHandler struct {
	billingClient *client.BillingClient
}

func NewAdminBillingAccountsHandler(billingClient *client.BillingClient) *AdminBillingAccountsHandler {
	_ = billingv1.BillingAccount{}
	return &AdminBillingAccountsHandler{
		billingClient: billingClient,
	}
}

type BillingAccountResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

func (h *AdminBillingAccountsHandler) GetBillingAccount(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	account, err := h.billingClient.GetBillingAccountByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	c.JSON(http.StatusOK, billingAccountToResponse(account))
}

func (h *AdminBillingAccountsHandler) CreateBillingAccount(c *gin.Context) {
	var req struct {
		UserID        string  `json:"user_id" binding:"required"`
		InitialCredit float64 `json:"initial_credit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	account, err := h.billingClient.CreateBillingAccount(c.Request.Context(), req.UserID, req.InitialCredit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, billingAccountToResponse(account))
}

func (h *AdminBillingAccountsHandler) AdjustBalance(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.billingClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing service unavailable"})
		return
	}

	account, err := h.billingClient.GetBillingAccountByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found, create one first"})
		return
	}

	newBalance := account.Balance + req.Amount
	if newBalance < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	updated, err := h.billingClient.UpdateBillingAccountBalance(c.Request.Context(), account.Id, newBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billingAccountToResponse(updated))
}

func billingAccountToResponse(a *billingv1.BillingAccount) BillingAccountResponse {
	return BillingAccountResponse{
		ID:       a.Id,
		UserID:   a.UserId,
		Balance:  a.Balance,
		Currency: a.Currency,
		Status:   a.Status,
	}
}
