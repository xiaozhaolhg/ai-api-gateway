package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// Gateway-level errors
	ErrProviderTimeout     ErrorCode = "gateway_timeout"
	ErrInvalidCredentials  ErrorCode = "invalid_api_key"
	ErrModelNotFound       ErrorCode = "model_not_found"
	ErrProviderUnavailable ErrorCode = "provider_error"
	ErrRateLimitExceeded   ErrorCode = "rate_limit_exceeded"
	ErrAuthDenied          ErrorCode = "insufficient_permissions"
	ErrInternal            ErrorCode = "internal_error"
	ErrBadRequest          ErrorCode = "bad_request"
	ErrServiceUnavailable  ErrorCode = "service_unavailable"
)

// GatewayError represents a structured gateway error
type GatewayError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

// Error implements the error interface
func (e *GatewayError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the cause error for error chaining
func (e *GatewayError) Unwrap() error {
	return e.Cause
}

// HTTPStatus returns the HTTP status code for this error
func (e *GatewayError) HTTPStatus() int {
	switch e.Code {
	case ErrInvalidCredentials:
		return http.StatusUnauthorized // 401
	case ErrAuthDenied:
		return http.StatusForbidden // 403
	case ErrModelNotFound:
		return http.StatusNotFound // 404
	case ErrRateLimitExceeded:
		return http.StatusTooManyRequests // 429
	case ErrProviderUnavailable:
		return http.StatusBadGateway // 502
	case ErrProviderTimeout:
		return http.StatusGatewayTimeout // 504
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable // 503
	case ErrBadRequest:
		return http.StatusBadRequest // 400
	default:
		return http.StatusInternalServerError // 500
	}
}

// Response returns the JSON error response
func (e *GatewayError) Response() map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"details": e.Details,
		},
	}
}

// ResponseJSON returns the JSON-encoded error response
func (e *GatewayError) ResponseJSON() []byte {
	data, _ := json.Marshal(e.Response())
	return data
}

// New creates a new GatewayError
func New(code ErrorCode, message string) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new GatewayError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap creates a new GatewayError wrapping an existing error
func Wrap(code ErrorCode, message string, cause error) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: cause.Error(),
	}
}

// Wrapf creates a new GatewayError wrapping an existing error with formatted message
func Wrapf(code ErrorCode, cause error, format string, args ...interface{}) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
		Details: cause.Error(),
	}
}

// FromGRPC converts a gRPC error to a GatewayError
func FromGRPC(err error) *GatewayError {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC status error
		return Wrap(ErrInternal, "internal error", err)
	}

	switch st.Code() {
	case codes.DeadlineExceeded:
		return Wrap(ErrProviderTimeout, "provider request timed out", err)
	case codes.Unauthenticated:
		return Wrap(ErrInvalidCredentials, "invalid or missing API key", err)
	case codes.NotFound:
		return Wrap(ErrModelNotFound, "model not found", err)
	case codes.Unavailable:
		return Wrap(ErrProviderUnavailable, "provider service unavailable", err)
	case codes.ResourceExhausted:
		return Wrap(ErrRateLimitExceeded, "rate limit exceeded", err)
	case codes.PermissionDenied:
		return Wrap(ErrAuthDenied, "insufficient permissions", err)
	case codes.InvalidArgument:
		return Wrap(ErrBadRequest, "invalid request", err)
	default:
		return Wrap(ErrInternal, "internal error", err)
	}
}

// FromGRPCWithContext converts a gRPC error with additional context
func FromGRPCWithContext(err error, context string) *GatewayError {
	if err == nil {
		return nil
	}

	gatewayErr := FromGRPC(err)
	if gatewayErr.Details != "" {
		gatewayErr.Details = fmt.Sprintf("%s: %s", context, gatewayErr.Details)
	} else {
		gatewayErr.Details = context
	}
	return gatewayErr
}

// Helper functions for common error types

// NewTimeoutError creates a timeout error
func NewTimeoutError(service string, duration string) *GatewayError {
	return Newf(ErrProviderTimeout, "%s request timed out after %s", service, duration)
}

// NewAuthError creates an authentication error
func NewAuthError(message string) *GatewayError {
	return New(ErrInvalidCredentials, message)
}

// NewModelNotFoundError creates a model not found error
func NewModelNotFoundError(model string) *GatewayError {
	return Newf(ErrModelNotFound, "model '%s' not found", model)
}

// NewProviderUnavailableError creates a provider unavailable error
func NewProviderUnavailableError(provider string) *GatewayError {
	return Newf(ErrProviderUnavailable, "provider '%s' is unavailable", provider)
}

// NewServiceUnavailableError creates a service unavailable error
func NewServiceUnavailableError(service string) *GatewayError {
	return Newf(ErrServiceUnavailable, "service '%s' is unavailable", service)
}
