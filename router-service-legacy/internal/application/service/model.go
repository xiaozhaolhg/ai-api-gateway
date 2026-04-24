package service

import (
	"github.com/ai-api-gateway/router-service-legacy/internal/domain/entity"
	"github.com/ai-api-gateway/router-service-legacy/internal/domain/port"
)

type ModelService struct {
	providers []port.Provider
}

func NewModelService(providers []port.Provider) *ModelService {
	return &ModelService{providers: providers}
}

func (s *ModelService) ListModels() ([]entity.Model, error) {
	var models []entity.Model
	for _, p := range s.providers {
		pModels, err := p.ListModels()
		if err != nil {
			continue
		}
		models = append(models, pModels...)
	}
	return models, nil
}