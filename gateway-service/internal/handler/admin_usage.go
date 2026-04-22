package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// AdminUsageHandler handles admin usage requests
type AdminUsageHandler struct {
	billingClient *client.BillingClient
}

// NewAdminUsageHandler creates a new admin usage handler
func NewAdminUsageHandler(billingClient *client.BillingClient) *AdminUsageHandler {
	return &AdminUsageHandler{
		billingClient: billingClient,
	}
}

// ServeHTTP handles HTTP requests for /admin/usage
func (h *AdminUsageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getUsage retrieves usage records
func (h *AdminUsageHandler) getUsage(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	resp, err := h.billingClient.GetUsage(r.Context(), userID, int32(page), int32(pageSize))
	if err != nil {
		http.Error(w, "Failed to get usage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
