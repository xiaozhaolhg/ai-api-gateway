package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
)

func TestCreateProviderServer(t *testing.T) {
	client, listener, server, err := CreateProviderServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer server.Stop()
	defer listener.Close()

	// Test that we can make a request
	resp, err := client.ForwardRequest(context.Background(), &providerv1.ForwardRequestRequest{
		ProviderId:  "ollama",
		RequestBody: []byte(`{"model":"llama2"}`),
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
