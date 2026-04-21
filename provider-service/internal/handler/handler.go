package handler

import (
	"context"
)

// Handler implements the gRPC ProviderService interface
// TODO: Implement this after buf generate creates the stubs
type Handler struct {
	// TODO: Add dependencies
}

// NewHandler creates a new Handler
func NewHandler() *Handler {
	return &Handler{}
}

// TODO: Implement gRPC methods after proto generation
// ForwardRequest, StreamRequest, CreateProvider, etc.

func (h *Handler) Shutdown(ctx context.Context) error {
	// Cleanup logic
	return nil
}
