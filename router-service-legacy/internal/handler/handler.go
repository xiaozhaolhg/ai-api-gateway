package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ai-api-gateway/router-service/internal/application/service"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/config"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/provider"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	chatService  *service.ChatCompletionService
	modelService *service.ModelService
	registry     *provider.ProviderRegistry
	config       *config.Config
}

func Setup(router *gin.Engine, cfg *config.Config) {
	// Initialize provider registry
	registry := provider.NewProviderRegistry()

	// Register built-in provider factories
	ollamaFactory := provider.NewOllamaFactory()
	if err := registry.Register(ollamaFactory); err != nil {
		log.Fatalf("Failed to register Ollama factory: %v", err)
	}

	opencodeFactory := provider.NewOpenCodeZenFactory()
	if err := registry.Register(opencodeFactory); err != nil {
		log.Fatalf("Failed to register OpenCode Zen factory: %v", err)
	}

	// Create providers from config using registry
	var providers []port.Provider
	enabledProviders := cfg.GetEnabledProviders()

	for providerType, settings := range enabledProviders {
		p, err := registry.Create(providerType, settings)
		if err != nil {
			log.Printf("Failed to create provider %s: %v", providerType, err)
			continue
		}
		log.Printf("Successfully initialized provider: %s", providerType)
		providers = append(providers, p)
	}

	if len(providers) == 0 {
		log.Println("Warning: No providers enabled")
	}

	routerService := service.NewModelRouter(providers)
	chatService := service.NewChatCompletionService(routerService)
	modelService := service.NewModelService(providers)

	h := &Handler{
		chatService:  chatService,
		modelService: modelService,
		registry:     registry,
		config:       cfg,
	}

	router.Use(corsMiddleware())
	router.GET("/health", h.healthHandler)
	router.POST("/v1/chat/completions", h.chatCompletionHandler)
	router.GET("/v1/models", h.modelsHandler)
	router.GET("/v1/providers", h.providersHandler)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) chatCompletionHandler(c *gin.Context) {
	var req struct {
		Model       string      `json:"model" binding:"required"`
		Messages    []struct {
			Role    string `json:"role" binding:"required"`
			Content string `json:"content" binding:"required"`
		} `json:"messages" binding:"required,min=1"`
		Stream      bool    `json:"stream"`
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"max_tokens"`
		TopP        float64 `json:"top_p"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages := make([]port.Message, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = port.Message{Role: m.Role, Content: m.Content}
	}

	chatReq := port.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}

	if req.Stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")

		ch, err := h.chatService.HandleStreamChatCompletion(c.Request.Context(), chatReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Stream(func(w io.Writer) bool {
			select {
			case chunk, ok := <-ch:
				if !ok {
					w.Write([]byte("data: [DONE]\n\n"))
					return false
				}
				if chunk.Done {
					w.Write([]byte("data: [DONE]\n\n"))
					return false
				}
				data, _ := json.Marshal(chunk)
				w.Write([]byte("data: " + string(data) + "\n\n"))
				return true
			case <-c.Request.Context().Done():
				return false
			}
		})
	} else {
		resp, err := h.chatService.HandleChatCompletion(c.Request.Context(), chatReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) modelsHandler(c *gin.Context) {
	models, err := h.modelService.ListModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]map[string]interface{}, len(models))
	for i, m := range models {
		data[i] = map[string]interface{}{
			"id":         m.ID,
			"object":     "model",
			"created":   0,
			"owned_by":   m.Provider,
			"provider":   m.Provider,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   data,
	})
}

func (h *Handler) providersHandler(c *gin.Context) {
	// Get all registered provider types
	types := h.registry.ListTypes()

	// Get defaults for all providers
	defaults := h.registry.GetDefaults()

	// Build response with provider info
	providers := make([]map[string]interface{}, 0, len(types))
	for _, providerType := range types {
		factory, err := h.registry.GetFactory(providerType)
		if err != nil {
			continue
		}

		providerInfo := map[string]interface{}{
			"type":        factory.Type(),
			"description": factory.Description(),
		}

		// Check if provider is configured
		if settings, ok := h.config.Provider.Providers[providerType]; ok {
			providerInfo["configured"] = true
			providerInfo["enabled"] = settings.Enabled
			providerInfo["endpoint"] = settings.Endpoint
			providerInfo["has_api_key"] = settings.APIKey != ""
		} else {
			providerInfo["configured"] = false
			providerInfo["enabled"] = false
		}

		// Add defaults
		if defaultSettings, ok := defaults[providerType]; ok {
			providerInfo["defaults"] = map[string]interface{}{
				"endpoint": defaultSettings.Endpoint,
				"enabled":   defaultSettings.Enabled,
			}
		}

		providers = append(providers, providerInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
	})
}