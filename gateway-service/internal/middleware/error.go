package middleware

import (
	stderrors "errors"

	"github.com/ai-api-gateway/gateway-service/internal/errors"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware handles error translation and formatting
type ErrorMiddleware struct{}

// NewErrorMiddleware creates a new error middleware
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// Middleware returns the middleware function
func (m *ErrorMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors in the context
		if len(c.Errors) > 0 {
			// Handle the first error
			err := c.Errors[0].Err
			m.handleError(c, err)
		}
	}
}

// handleError translates and formats errors
func (m *ErrorMiddleware) handleError(c *gin.Context, err error) {
	var gatewayErr *errors.GatewayError

	// Check if it's already a GatewayError
	if ge, ok := err.(*errors.GatewayError); ok {
		gatewayErr = ge
	} else {
		// Check if it wraps a GatewayError
		if ge, ok := stderrors.Unwrap(err).(*errors.GatewayError); ok {
			gatewayErr = ge
		} else {
			// Convert to internal error
			gatewayErr = errors.Wrap(errors.ErrInternal, "internal error", err)
		}
	}

	// Send error response
	c.JSON(gatewayErr.HTTPStatus(), gatewayErr.Response())
}

// ErrorResponse writes a gateway error response directly
func ErrorResponse(c *gin.Context, err *errors.GatewayError) {
	c.JSON(err.HTTPStatus(), err.Response())
}

// AbortWithError aborts the request with a gateway error
func AbortWithError(c *gin.Context, err *errors.GatewayError) {
	c.AbortWithStatusJSON(err.HTTPStatus(), err.Response())
}

// HandleGRPCError converts a gRPC error to HTTP response
func HandleGRPCError(c *gin.Context, err error, context string) {
	if err == nil {
		return
	}

	gatewayErr := errors.FromGRPCWithContext(err, context)
	AbortWithError(c, gatewayErr)
}
