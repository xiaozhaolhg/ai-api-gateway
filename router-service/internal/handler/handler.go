package handler

import (
	"encoding/json"
	"io"
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
}

func Setup(router *gin.Engine, cfg *config.Config) {
	var providers []port.Provider

	for _, p := range cfg.GetEnabledProviders() {
		if p.Endpoint == "" {
			continue
		}
		if p.Endpoint == "http://localhost:11434" || p.Endpoint == "http://host.docker.internal:11434" || p.Endpoint == "http://172.17.0.1:11434" {
			providers = append(providers, provider.NewOllamaProvider(p))
		} else {
			providers = append(providers, provider.NewOpenCodeZenProvider(p))
		}
	}

	routerService := service.NewModelRouter(providers)
	chatService := service.NewChatCompletionService(routerService)
	modelService := service.NewModelService(providers)

	h := &Handler{
		chatService:  chatService,
		modelService: modelService,
	}

	router.Use(corsMiddleware())
	router.GET("/health", h.healthHandler)
	router.POST("/v1/chat/completions", h.chatCompletionHandler)
	router.GET("/v1/models", h.modelsHandler)
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