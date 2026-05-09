package repository

import (
	"log"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
	"gorm.io/gorm"
)

// ProviderRepository implements the ProviderRepository interface using GORM
type ProviderRepository struct {
	db *gorm.DB
}

// NewProviderRepository creates a new ProviderRepository
func NewProviderRepository(db *gorm.DB) port.ProviderRepository {
	return &ProviderRepository{db: db}
}

// Create creates a new provider
func (r *ProviderRepository) Create(provider *entity.Provider) error {
	return r.db.Create(provider).Error
}

// GetByID retrieves a provider by ID
func (r *ProviderRepository) GetByID(id string) (*entity.Provider, error) {
	var provider entity.Provider
	err := r.db.Where("id = ?", id).First(&provider).Error
	if err != nil {
		// Try by name as fallback
		err2 := r.db.Where("name = ?", id).First(&provider).Error
		if err2 != nil {
			return nil, err
		}
	}
	return &provider, nil
}

// GetByType retrieves a provider by type
func (r *ProviderRepository) GetByType(providerType string) (*entity.Provider, error) {
	var provider entity.Provider
	err := r.db.Where("type = ?", providerType).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetByName retrieves a provider by name
func (r *ProviderRepository) GetByName(name string) (*entity.Provider, error) {
	var provider entity.Provider
	err := r.db.Where("name = ?", name).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// Update updates a provider
func (r *ProviderRepository) Update(provider *entity.Provider) error {
	return r.db.Save(provider).Error
}

// Delete deletes a provider by ID
func (r *ProviderRepository) Delete(id string) error {
	return r.db.Delete(&entity.Provider{}, "id = ?", id).Error
}

// List retrieves providers with pagination
func (r *ProviderRepository) List(page, pageSize int) ([]*entity.Provider, int, error) {
	var providers []*entity.Provider
	var total int64

	if err := r.db.Model(&entity.Provider{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&providers).Error
	if err != nil {
		return nil, 0, err
	}

	return providers, int(total), nil
}

func (r *ProviderRepository) FindByModel(model string) ([]*entity.Provider, error) {
	var allProviders []*entity.Provider
	if err := r.db.Find(&allProviders).Error; err != nil {
		return nil, err
	}

	log.Printf("[FindByModel] Searching for model=%s, total providers=%d", model, len(allProviders))
	for _, p := range allProviders {
		log.Printf("[FindByModel] Provider ID=%s, Name=%s, Models=%v", p.ID, p.Name, p.Models)
	}

	var matched []*entity.Provider
	for _, p := range allProviders {
		if len(p.Models) == 0 {
			continue
		}
		for _, m := range p.Models {
			log.Printf("[FindByModel] Comparing model=%s with provider model=%s", model, m)
			if m == model {
				log.Printf("[FindByModel] MATCH found: provider ID=%s", p.ID)
				matched = append(matched, p)
				break
			}
		}
	}
	return matched, nil
}
