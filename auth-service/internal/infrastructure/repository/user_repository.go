package repository

import (
	"log"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

// UserRepository implements the UserRepository interface using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	log.Printf("[DEBUG] GetByUsername query: username=%s", username)
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Printf("[DEBUG] GetByUsername result: error=%v", err)
		return nil, err
	}
	log.Printf("[DEBUG] GetByUsername found: user_id=%s, username=%s", user.ID, user.Username)
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}

// List retrieves users with pagination
func (r *UserRepository) List(page, pageSize int) ([]*entity.User, int, error) {
	var users []*entity.User
	var total int64

	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}
