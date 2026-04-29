package port

import (
	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
)

// ProviderRepository defines the interface for provider persistence operations
type ProviderRepository interface {
	Create(provider *entity.Provider) error
	GetByID(id string) (*entity.Provider, error)
	GetByName(name string) (*entity.Provider, error)
	GetByType(providerType string) (*entity.Provider, error)
	Update(provider *entity.Provider) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.Provider, int, error)
}
