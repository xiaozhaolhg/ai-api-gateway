package repository

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

// APIKeyRepository implements the APIKeyRepository interface using GORM
type APIKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new APIKeyRepository
func NewAPIKeyRepository(db *gorm.DB) port.APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create creates a new API key
func (r *APIKeyRepository) Create(apiKey *entity.APIKey) error {
	return r.db.Create(apiKey).Error
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(id string) (*entity.APIKey, error) {
	var apiKey entity.APIKey
	err := r.db.Where("id = ?", id).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetByKeyHash retrieves an API key by its hash
func (r *APIKeyRepository) GetByKeyHash(keyHash string) (*entity.APIKey, error) {
	var apiKey entity.APIKey
	err := r.db.Where("key_hash = ?", keyHash).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetByUserID retrieves API keys for a user with pagination
func (r *APIKeyRepository) GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error) {
	var apiKeys []*entity.APIKey
	var total int64

	if err := r.db.Model(&entity.APIKey{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&apiKeys).Error
	if err != nil {
		return nil, 0, err
	}

	return apiKeys, int(total), nil
}

// Delete deletes an API key by ID
func (r *APIKeyRepository) Delete(id string) error {
	return r.db.Delete(&entity.APIKey{}, "id = ?", id).Error
}
