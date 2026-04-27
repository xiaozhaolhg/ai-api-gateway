package port

import (
	"testing"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
)

// MockAdapter is a mock implementation of ProviderAdapter for testing
type MockAdapter struct {
	TransformRequestFunc  func(request []byte, headers map[string]string) ([]byte, map[string]string, error)
	TransformResponseFunc func(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error)
	CountTokensFunc       func(request []byte, response []byte, isStreaming bool) (int64, int64, error)
}

func (m *MockAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
	if m.TransformRequestFunc != nil {
		return m.TransformRequestFunc(request, headers)
	}
	return request, headers, nil
}

func (m *MockAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if m.TransformResponseFunc != nil {
		return m.TransformResponseFunc(response, isStreaming, accumulatedTokens)
	}
	return response, accumulatedTokens, !isStreaming, nil
}

func (m *MockAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if m.CountTokensFunc != nil {
		return m.CountTokensFunc(request, response, isStreaming)
	}
	return 0, 0, nil
}

func TestProviderAdapter_TransformRequest(t *testing.T) {
	mock := &MockAdapter{
		TransformRequestFunc: func(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
			// Simple transformation: add a prefix
			transformed := append([]byte("transformed:"), request...)
			return transformed, headers, nil
		},
	}

	request := []byte("test request")
	headers := map[string]string{"Content-Type": "application/json"}

	transformed, newHeaders, err := mock.TransformRequest(request, headers)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	if string(transformed) != "transformed:test request" {
		t.Errorf("Expected transformed request, got %s", string(transformed))
	}

	if newHeaders["Content-Type"] != "application/json" {
		t.Error("Expected headers to be preserved")
	}
}

func TestProviderAdapter_TransformResponse(t *testing.T) {
	mock := &MockAdapter{
		TransformResponseFunc: func(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
			// Simple transformation: add a prefix
			transformed := append([]byte("transformed:"), response...)
			return transformed, accumulatedTokens, !isStreaming, nil
		},
	}

	response := []byte("test response")

	transformed, tokenCounts, isFinal, err := mock.TransformResponse(response, false, entity.TokenCounts{})
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	if string(transformed) != "transformed:test response" {
		t.Errorf("Expected transformed response, got %s", string(transformed))
	}

	if !isFinal {
		t.Error("Expected isFinal to be true for non-streaming")
	}

	if tokenCounts.Total() != 0 {
		t.Errorf("Expected empty token counts, got %d", tokenCounts.Total())
	}
}

func TestProviderAdapter_CountTokens(t *testing.T) {
	mock := &MockAdapter{
		CountTokensFunc: func(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
			// Simple counting: return length of request and response
			return int64(len(request)), int64(len(response)), nil
		},
	}

	request := []byte("test request")
	response := []byte("test response")

	reqTokens, respTokens, err := mock.CountTokens(request, response, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	if reqTokens != int64(len(request)) {
		t.Errorf("Expected request tokens %d, got %d", len(request), reqTokens)
	}

	if respTokens != int64(len(response)) {
		t.Errorf("Expected response tokens %d, got %d", len(response), respTokens)
	}
}

func TestProviderAdapter_CountTokens_Streaming(t *testing.T) {
	mock := &MockAdapter{
		CountTokensFunc: func(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
			// For streaming intermediate chunks, return 0, 0
			if isStreaming {
				return 0, 0, nil
			}
			return int64(len(request)), int64(len(response)), nil
		},
	}

	request := []byte("test request")
	response := []byte("test response")

	// Test streaming mode
	reqTokens, respTokens, err := mock.CountTokens(request, response, true)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	if reqTokens != 0 {
		t.Errorf("Expected 0 prompt tokens for streaming, got %d", reqTokens)
	}

	if respTokens != 0 {
		t.Errorf("Expected 0 completion tokens for streaming, got %d", respTokens)
	}
}

func TestProviderAdapter_ErrorHandling(t *testing.T) {
	t.Run("TransformRequest error", func(t *testing.T) {
		mock := &MockAdapter{
			TransformRequestFunc: func(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
				return nil, nil, &testError{"transform error"}
			},
		}

		_, _, err := mock.TransformRequest([]byte("test"), nil)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("TransformResponse error", func(t *testing.T) {
		mock := &MockAdapter{
			TransformResponseFunc: func(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
				return nil, entity.TokenCounts{}, false, &testError{"transform error"}
			},
		}

		_, _, _, err := mock.TransformResponse([]byte("test"), false, entity.TokenCounts{})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("CountTokens error", func(t *testing.T) {
		mock := &MockAdapter{
			CountTokensFunc: func(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
				return 0, 0, &testError{"count error"}
			},
		}

		_, _, err := mock.CountTokens([]byte("test"), []byte("test"), false)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
