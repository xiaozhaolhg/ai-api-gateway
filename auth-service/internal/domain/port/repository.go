package port

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.User, int, error)
}

// APIKeyRepository defines the interface for API key persistence operations
type APIKeyRepository interface {
	Create(apiKey *entity.APIKey) error
	GetByID(id string) (*entity.APIKey, error)
	GetByKeyHash(keyHash string) (*entity.APIKey, error)
	GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error)
	Delete(id string) error
}
