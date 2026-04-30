package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	err := New(ErrInvalidCredentials, "invalid api key")

	if err.Code != ErrInvalidCredentials {
		t.Errorf("expected code %s, got %s", ErrInvalidCredentials, err.Code)
	}
	if err.Message != "invalid api key" {
		t.Errorf("expected message 'invalid api key', got %s", err.Message)
	}
}

func TestNewf(t *testing.T) {
	err := Newf(ErrModelNotFound, "model '%s' not found", "gpt-4")

	if err.Code != ErrModelNotFound {
		t.Errorf("expected code %s, got %s", ErrModelNotFound, err.Code)
	}
	if err.Message != "model 'gpt-4' not found" {
		t.Errorf("expected formatted message, got %s", err.Message)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := Wrap(ErrServiceUnavailable, "auth service unavailable", cause)

	if err.Code != ErrServiceUnavailable {
		t.Errorf("expected code %s, got %s", ErrServiceUnavailable, err.Code)
	}
	if err.Cause != cause {
		t.Error("expected cause to be set")
	}
	if err.Details != "connection refused" {
		t.Errorf("expected details 'connection refused', got %s", err.Details)
	}
}

func TestGatewayError_Error(t *testing.T) {
	err := New(ErrInvalidCredentials, "invalid api key")
	expected := "invalid_api_key: invalid api key"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}

	// With details
	err2 := Wrap(ErrInternal, "internal error", errors.New("db connection failed"))
	if err2.Error() == "" {
		t.Error("expected non-empty error string")
	}
}

func TestGatewayError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := Wrap(ErrInternal, "wrapped", cause)

	if err.Unwrap() != cause {
		t.Error("expected Unwrap to return cause")
	}
}

func TestGatewayError_HTTPStatus(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected int
	}{
		{ErrInvalidCredentials, http.StatusUnauthorized},
		{ErrAuthDenied, http.StatusForbidden},
		{ErrModelNotFound, http.StatusNotFound},
		{ErrRateLimitExceeded, http.StatusTooManyRequests},
		{ErrProviderUnavailable, http.StatusBadGateway},
		{ErrProviderTimeout, http.StatusGatewayTimeout},
		{ErrServiceUnavailable, http.StatusServiceUnavailable},
		{ErrBadRequest, http.StatusBadRequest},
		{ErrInternal, http.StatusInternalServerError},
		{"unknown_code", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		err := New(tt.code, "test")
		if err.HTTPStatus() != tt.expected {
			t.Errorf("code %s: expected status %d, got %d", tt.code, tt.expected, err.HTTPStatus())
		}
	}
}

func TestGatewayError_Response(t *testing.T) {
	err := New(ErrInvalidCredentials, "invalid api key")
	resp := err.Response()

	errMap, ok := resp["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error map in response")
	}

	code, ok := errMap["code"].(ErrorCode)
	if !ok {
		// Try as string if not ErrorCode
		codeStr, _ := errMap["code"].(string)
		code = ErrorCode(codeStr)
	}
	if code != ErrInvalidCredentials {
		t.Errorf("expected code %s, got %v", ErrInvalidCredentials, errMap["code"])
	}
}

func TestGatewayError_ResponseJSON(t *testing.T) {
	err := New(ErrInvalidCredentials, "invalid api key")
	jsonData := err.ResponseJSON()

	if len(jsonData) == 0 {
		t.Error("expected non-empty JSON")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
}

func TestFromGRPC(t *testing.T) {
	tests := []struct {
		name         string
		grpcErr      error
		expectedCode ErrorCode
	}{
		{
			name:         "nil error",
			grpcErr:      nil,
			expectedCode: "",
		},
		{
			name:         "deadline exceeded",
			grpcErr:      status.Error(codes.DeadlineExceeded, "timeout"),
			expectedCode: ErrProviderTimeout,
		},
		{
			name:         "unauthenticated",
			grpcErr:      status.Error(codes.Unauthenticated, "auth failed"),
			expectedCode: ErrInvalidCredentials,
		},
		{
			name:         "not found",
			grpcErr:      status.Error(codes.NotFound, "model not found"),
			expectedCode: ErrModelNotFound,
		},
		{
			name:         "unavailable",
			grpcErr:      status.Error(codes.Unavailable, "service down"),
			expectedCode: ErrProviderUnavailable,
		},
		{
			name:         "resource exhausted",
			grpcErr:      status.Error(codes.ResourceExhausted, "rate limited"),
			expectedCode: ErrRateLimitExceeded,
		},
		{
			name:         "permission denied",
			grpcErr:      status.Error(codes.PermissionDenied, "no access"),
			expectedCode: ErrAuthDenied,
		},
		{
			name:         "invalid argument",
			grpcErr:      status.Error(codes.InvalidArgument, "bad request"),
			expectedCode: ErrBadRequest,
		},
		{
			name:         "internal error",
			grpcErr:      status.Error(codes.Internal, "server error"),
			expectedCode: ErrInternal,
		},
		{
			name:         "non-status error",
			grpcErr:      errors.New("random error"),
			expectedCode: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPC(tt.grpcErr)
			if tt.grpcErr == nil {
				if result != nil {
					t.Error("expected nil for nil error")
				}
				return
			}
			if result.Code != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, result.Code)
			}
		})
	}
}

func TestFromGRPCWithContext(t *testing.T) {
	grpcErr := status.Error(codes.Unavailable, "connection refused")
	result := FromGRPCWithContext(grpcErr, "auth service")

	if result.Code != ErrProviderUnavailable {
		t.Errorf("expected code %s, got %s", ErrProviderUnavailable, result.Code)
	}
	if result.Details == "" {
		t.Error("expected details with context")
	}
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		err      *GatewayError
		expected ErrorCode
	}{
		{
			name:     "NewTimeoutError",
			err:      NewTimeoutError("auth", "5s"),
			expected: ErrProviderTimeout,
		},
		{
			name:     "NewAuthError",
			err:      NewAuthError("login failed"),
			expected: ErrInvalidCredentials,
		},
		{
			name:     "NewModelNotFoundError",
			err:      NewModelNotFoundError("gpt-4"),
			expected: ErrModelNotFound,
		},
		{
			name:     "NewProviderUnavailableError",
			err:      NewProviderUnavailableError("ollama"),
			expected: ErrProviderUnavailable,
		},
		{
			name:     "NewServiceUnavailableError",
			err:      NewServiceUnavailableError("billing"),
			expected: ErrServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expected {
				t.Errorf("expected code %s, got %s", tt.expected, tt.err.Code)
			}
		})
	}
}

