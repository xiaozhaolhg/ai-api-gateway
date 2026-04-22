package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// ModelsHandler handles model listing requests
type ModelsHandler struct {
	providerClient *client.ProviderClient
}

// NewModelsHandler creates a new models handler
func NewModelsHandler(providerClient *client.ProviderClient) *ModelsHandler {
	return &ModelsHandler{
		providerClient: providerClient,
	}
}

// ServeHTTP handles HTTP requests for /v1/models
func (h *ModelsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// List all providers to get available models
	providersResp, err := h.providerClient.ListProviders(r.Context(), 1, 100)
	if err != nil {
		http.Error(w, "Failed to list providers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Collect all models from all providers
	var models []map[string]interface{}
	for _, provider := range providersResp.Providers {
		for _, model := range provider.Models {
			models = append(models, map[string]interface{}{
				"id":      model,
				"object":  "model",
				"created": provider.CreatedAt,
				"owned_by": provider.Name,
			})
		}
	}

	// Return OpenAI-compatible response
	response := map[string]interface{}{
		"object": "list",
		"data":   models,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
