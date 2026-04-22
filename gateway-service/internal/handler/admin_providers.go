package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-api-gateway/api/gen/provider/v1"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// AdminProvidersHandler handles admin provider management requests
type AdminProvidersHandler struct {
	providerClient *client.ProviderClient
}

// NewAdminProvidersHandler creates a new admin providers handler
func NewAdminProvidersHandler(providerClient *client.ProviderClient) *AdminProvidersHandler {
	return &AdminProvidersHandler{
		providerClient: providerClient,
	}
}

// ServeHTTP handles HTTP requests for /admin/providers
func (h *AdminProvidersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProviders(w, r)
	case http.MethodPost:
		h.createProvider(w, r)
	case http.MethodPut:
		h.updateProvider(w, r)
	case http.MethodDelete:
		h.deleteProvider(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listProviders lists all providers
func (h *AdminProvidersHandler) listProviders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	resp, err := h.providerClient.ListProviders(r.Context(), int32(page), int32(pageSize))
	if err != nil {
		http.Error(w, "Failed to list providers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createProvider creates a new provider
func (h *AdminProvidersHandler) createProvider(w http.ResponseWriter, r *http.Request) {
	var provider providerv1.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.providerClient.CreateProvider(r.Context(), &provider)
	if err != nil {
		http.Error(w, "Failed to create provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// updateProvider updates an existing provider
func (h *AdminProvidersHandler) updateProvider(w http.ResponseWriter, r *http.Request) {
	var provider providerv1.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.providerClient.UpdateProvider(r.Context(), &provider)
	if err != nil {
		http.Error(w, "Failed to update provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// deleteProvider deletes a provider
func (h *AdminProvidersHandler) deleteProvider(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing provider id", http.StatusBadRequest)
		return
	}

	err := h.providerClient.DeleteProvider(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
