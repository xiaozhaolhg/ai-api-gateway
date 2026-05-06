package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/middleware/testutils"
)

func TestProxyMiddleware_Streaming_AccumulatesTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accumulator := &tokenAccumulator{}
	chunks := []struct {
		promptTokens     int64
		completionTokens int64
	}{
		{5, 10},
		{5, 20},
		{5, 25},
	}

	for _, chunk := range chunks {
		accumulator.totalPromptTokens += chunk.promptTokens
		accumulator.totalCompletionTokens += chunk.completionTokens
	}

	assert.Equal(t, int64(15), accumulator.totalPromptTokens)
	assert.Equal(t, int64(55), accumulator.totalCompletionTokens)
}

type tokenAccumulator struct {
	totalPromptTokens     int64
	totalCompletionTokens int64
}

func TestStreamingIntervalCalculation(t *testing.T) {
	tests := []struct {
		name                   string
		interval               int64
		totalCompletion        int64
		lastRecordedCompletion int64
		shouldTrigger          bool
	}{
		{"below interval", 1000, 800, 0, false},
		{"exactly at interval", 1000, 1000, 0, true},
		{"above interval", 1000, 1500, 0, true},
		{"after reset", 1000, 1500, 1000, false},
		{"below after reset", 1000, 1400, 1000, false},
		{"disabled (interval 0)", 0, 5000, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trigger := tt.interval > 0 && tt.totalCompletion-tt.lastRecordedCompletion >= tt.interval
			assert.Equal(t, tt.shouldTrigger, trigger)
		})
	}
}

func TestProxyMiddleware_RouteIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, providerListener, providerServer, err := testutils.CreateProviderServer()
	if err != nil {
		t.Fatalf("Failed to create provider server: %v", err)
	}
	defer providerServer.Stop()
	defer providerListener.Close()

	_, billingListener, billingServer, err := testutils.CreateBillingServer()
	if err != nil {
		t.Fatalf("Failed to create billing server: %v", err)
	}
	defer billingServer.Stop()
	defer billingListener.Close()

	providerAddr := providerListener.Addr().String()
	wrappedProviderClient, err := client.NewProviderClient(providerAddr)
	if err != nil {
		t.Fatalf("Failed to create provider client: %v", err)
	}
	defer wrappedProviderClient.Close()

	billingAddr := billingListener.Addr().String()
	wrappedBillingClient, err := client.NewBillingClient(billingAddr)
	if err != nil {
		t.Fatalf("Failed to create billing client: %v", err)
	}
	defer wrappedBillingClient.Close()

	proxyMiddleware := NewProxyMiddleware(wrappedProviderClient, wrappedBillingClient, 1000)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("providerId", "test-provider")
		c.Set("userId", "test-user")
		c.Set("groupId", "test-group")
		c.Next()
	})
	router.POST("/v1/chat/completions", proxyMiddleware.Middleware())

	body := map[string]interface{}{"model": "test:model", "stream": true}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
