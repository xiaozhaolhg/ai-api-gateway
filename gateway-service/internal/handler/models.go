package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// Model represents an AI model in the OpenAI-compatible format
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsListResp represents the OpenAI-compatible models list response
type ModelsListResp struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// cachedModels holds cached model data with timestamp
type cachedModels struct {
	models    []Model
	cachedAt  time.Time
}

// ModelsHandler handles model listing with aggregation from providers
type ModelsHandler struct {
	providerClient *client.ProviderClient
	cache          *cachedModels
	cacheMutex     sync.RWMutex
	cacheTTL       time.Duration
}

// NewModelsHandler creates a new models handler with caching
func NewModelsHandler(providerClient *client.ProviderClient) *ModelsHandler {
	return &ModelsHandler{
		providerClient: providerClient,
		cacheTTL:       5 * time.Minute, // 5 minute TTL as per spec
	}
}

// ListModels returns an aggregated list of all available models from providers
func (h *ModelsHandler) ListModels(c *gin.Context) {
	// Check cache first
	if cached := h.getCachedModels(); cached != nil {
		c.JSON(http.StatusOK, ModelsListResp{
			Object: "list",
			Data:   cached,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Fetch providers with their models
	providersResp, err := h.providerClient.ListProviders(ctx, 1, 100)
	if err != nil {
		// Return empty list on error (graceful fallback)
		c.JSON(http.StatusOK, ModelsListResp{
			Object: "list",
			Data:   []Model{},
		})
		return
	}

	// Aggregate models from all providers
	var models []Model
	now := time.Now().Unix()

	for _, provider := range providersResp.Providers {
		for _, modelID := range provider.Models {
			// Create model entry with provider prefix (e.g., "ollama:llama2")
			fullModelID := modelID
			if provider.Type != "" {
				fullModelID = provider.Type + ":" + modelID
			}

			models = append(models, Model{
				ID:      fullModelID,
				Object:  "model",
				Created: now,
				OwnedBy: provider.Name,
			})
		}
	}

	// Cache the results
	h.setCachedModels(models)

	c.JSON(http.StatusOK, ModelsListResp{
		Object: "list",
		Data:   models,
	})
}

// getCachedModels returns cached models if still valid
func (h *ModelsHandler) getCachedModels() []Model {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()

	if h.cache == nil {
		return nil
	}

	if time.Since(h.cache.cachedAt) > h.cacheTTL {
		return nil // Cache expired
	}

	return h.cache.models
}

// setCachedModels stores models in cache
func (h *ModelsHandler) setCachedModels(models []Model) {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()

	h.cache = &cachedModels{
		models:   models,
		cachedAt: time.Now(),
	}
}