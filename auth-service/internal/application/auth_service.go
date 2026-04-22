package application

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
)

// AuthService provides authentication and authorization logic
type AuthService struct {
	userRepo    port.UserRepository
	apiKeyRepo port.APIKeyRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo port.UserRepository, apiKeyRepo port.APIKeyRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

// GenerateAPIKey generates a random API key
func (s *AuthService) GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("sk-%x", b), nil
}

// HashAPIKey hashes an API key using SHA-256
func (s *AuthService) HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// ValidateAPIKey validates an API key and returns the user identity
func (s *AuthService) ValidateAPIKey(apiKey string) (*entity.User, []string, error) {
	keyHash := s.HashAPIKey(apiKey)
	apiKeyRecord, err := s.apiKeyRepo.GetByKeyHash(keyHash)
	if err != nil {
		return nil, nil, err
	}

	// Check if API key is expired
	if apiKeyRecord.ExpiresAt != nil && apiKeyRecord.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("API key expired")
	}

	user, err := s.userRepo.GetByID(apiKeyRecord.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user.Status != "active" {
		return nil, nil, fmt.Errorf("user is disabled")
	}

	return user, apiKeyRecord.Scopes, nil
}

// CheckModelAuthorization checks if a user is authorized to access a model
func (s *AuthService) CheckModelAuthorization(userID string, groupIDs []string, model string) (bool, []string, string) {
	// MVP: All active users are authorized to access all models
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, nil, "user not found"
	}

	if user.Status != "active" {
		return false, nil, "user disabled"
	}

	// MVP: Return all models as authorized
	// In Phase 2+, this would check permissions and group memberships
	return true, []string{"*"}, ""
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(name, email, role string) (*entity.User, error) {
	user := &entity.User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		Role:      role,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateAPIKey creates a new API key for a user
func (s *AuthService) CreateAPIKey(userID, name string) (string, string, error) {
	apiKey, err := s.GenerateAPIKey()
	if err != nil {
		return "", "", err
	}

	apiKeyRecord := &entity.APIKey{
		ID:        generateID(),
		UserID:    userID,
		KeyHash:   s.HashAPIKey(apiKey),
		Name:      name,
		Scopes:    []string{"read", "write"},
		CreatedAt: time.Now(),
	}

	if err := s.apiKeyRepo.Create(apiKeyRecord); err != nil {
		return "", "", err
	}

	return apiKeyRecord.ID, apiKey, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
