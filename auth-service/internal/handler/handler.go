package handler

import (
	"context"

	"github.com/ai-api-gateway/auth-service/internal/application"
)

// Handler implements the gRPC AuthService interface
// TODO: Implement this after buf generate creates the stubs
type Handler struct {
	authService *application.AuthService
}

// NewHandler creates a new Handler
func NewHandler(authService *application.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// TODO: Implement gRPC methods after proto generation
// ValidateAPIKey, CheckModelAuthorization, GetUser, CreateUser, etc.

func (h *Handler) Shutdown(ctx context.Context) error {
	// Cleanup logic
	return nil
}
