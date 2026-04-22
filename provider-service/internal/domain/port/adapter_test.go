package port

import (
	"testing"
)

// MockAdapter is a mock implementation of ProviderAdapter for testing
type MockAdapter struct {
	TransformRequestFunc    func(request []byte, headers map[string]string) ([]byte, map[string]string, error)
	TransformResponseFunc   func(response []byte) ([]byte, error)
	CountTokensFunc         func(request []byte, response []byte) (int64, int64, error)
}

func (m *MockAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
	if m.TransformRequestFunc != nil {
		return m.TransformRequestFunc(request, headers)
	}
	return request, headers, nil
}

func (m *MockAdapter) TransformResponse(response []byte) ([]byte, error) {
	if m.TransformResponseFunc != nil {
		return m.TransformResponseFunc(response)
	}
	return response, nil
}

func (m *MockAdapter) CountTokens(request []byte, response []byte) (int64, int64, error) {
	if m.CountTokensFunc != nil {
		return m.CountTokensFunc(request, response)
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
		TransformResponseFunc: func(response []byte) ([]byte, error) {
			// Simple transformation: add a prefix
			transformed := append([]byte("transformed:"), response...)
			return transformed, nil
		},
	}

	response := []byte("test response")

	transformed, err := mock.TransformResponse(response)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	if string(transformed) != "transformed:test response" {
		t.Errorf("Expected transformed response, got %s", string(transformed))
	}
}

func TestProviderAdapter_CountTokens(t *testing.T) {
	mock := &MockAdapter{
		CountTokensFunc: func(request []byte, response []byte) (int64, int64, error) {
			// Simple counting: return length of request and response
			return int64(len(request)), int64(len(response)), nil
		},
	}

	request := []byte("test request")
	response := []byte("test response")

	reqTokens, respTokens, err := mock.CountTokens(request, response)
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
			TransformResponseFunc: func(response []byte) ([]byte, error) {
				return nil, &testError{"transform error"}
			},
		}

		_, err := mock.TransformResponse([]byte("test"))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("CountTokens error", func(t *testing.T) {
		mock := &MockAdapter{
			CountTokensFunc: func(request []byte, response []byte) (int64, int64, error) {
				return 0, 0, &testError{"count error"}
			},
		}

		_, _, err := mock.CountTokens([]byte("test"), []byte("test"))
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
