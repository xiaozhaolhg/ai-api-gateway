package service

import (
	"fmt"

	"github.com/ai-api-gateway/router-service-legacy/internal/domain/port"
)

type ModelRouter struct {
	providers []port.Provider
}

func NewModelRouter(providers []port.Provider) *ModelRouter {
	return &ModelRouter{providers: providers}
}

func (r *ModelRouter) SelectProvider(model string) (port.Provider, error) {
	for _, p := range r.providers {
		if isModelForProvider(model, p.Name()) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no provider found for model: %s", model)
}

func isModelForProvider(model, provider string) bool {
	prefix := provider + ":"
	return len(model) > len(prefix) && model[:len(prefix)] == prefix
}