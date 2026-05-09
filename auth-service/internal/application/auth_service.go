package application

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"github.com/ai-api-gateway/pkg/cache"
)

func matchGlob(pattern, text string) bool {
	if pattern == "*" {
		return true
	}
	patternParts := strings.Split(pattern, ":")
	textParts := strings.Split(text, ":")
	if len(patternParts) != 2 || len(textParts) != 2 {
		return pattern == text
	}
	if patternParts[0] != "*" && patternParts[0] != textParts[0] {
		return false
	}
	modelPattern := patternParts[1]
	modelText := textParts[1]
	if strings.HasSuffix(modelPattern, "*") {
		prefix := strings.TrimSuffix(modelPattern, "*")
		return strings.HasPrefix(modelText, prefix)
	}
	return modelPattern == modelText
}

// AuthService provides authentication and authorization logic
type AuthService struct {
	userRepo      port.UserRepository
	apiKeyRepo    port.APIKeyRepository
	userGroupRepo port.UserGroupRepository
	tierRepo      port.TierRepository
	groupRepo     port.GroupRepository
	apiKeyCache   *cache.Cache[string, *entity.APIKey]
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo port.UserRepository, apiKeyRepo port.APIKeyRepository, userGroupRepo port.UserGroupRepository, tierRepo port.TierRepository, groupRepo port.GroupRepository) *AuthService {
	// Cache API keys for 5 minutes
	apiKeyCache := cache.New[string, *entity.APIKey](5 * time.Minute)
	apiKeyCache.StartCleanup(1 * time.Minute)

	return &AuthService{
		userRepo:      userRepo,
		apiKeyRepo:    apiKeyRepo,
		userGroupRepo: userGroupRepo,
		tierRepo:      tierRepo,
		groupRepo:     groupRepo,
		apiKeyCache:   apiKeyCache,
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
func (s *AuthService) ValidateAPIKey(apiKey string) (*entity.User, []string, []string, error) {
	keyHash := s.HashAPIKey(apiKey)

	// Try cache first
	apiKeyRecord, found := s.apiKeyCache.Get(keyHash)
	if !found {
		// Cache miss - fetch from repository
		var err error
		apiKeyRecord, err = s.apiKeyRepo.GetByKeyHash(keyHash)
		if err != nil {
			return nil, nil, nil, err
		}
		if apiKeyRecord == nil {
			return nil, nil, nil, fmt.Errorf("API key not found")
		}
		// Cache the result
		s.apiKeyCache.Set(keyHash, apiKeyRecord)
	}

	// Check if API key is expired
	if apiKeyRecord.ExpiresAt != nil && apiKeyRecord.ExpiresAt.Before(time.Now()) {
		return nil, nil, nil, fmt.Errorf("API key expired")
	}

	user, err := s.userRepo.GetByID(apiKeyRecord.UserID)
	if err != nil {
		return nil, nil, nil, err
	}

	if user.Status != "active" {
		return nil, nil, nil, fmt.Errorf("user is disabled")
	}

	// Populate group IDs from UserGroupMembership
	var groupIDs []string
	if s.userGroupRepo != nil {
		memberships, err := s.userGroupRepo.GetByUserID(user.ID)
		if err == nil && len(memberships) > 0 {
			groupIDs = make([]string, len(memberships))
			for i, m := range memberships {
				groupIDs[i] = m.GroupID
			}
		}
	}
	if groupIDs == nil {
		groupIDs = []string{}
	}

	return user, apiKeyRecord.Scopes, groupIDs, nil
}

// CheckModelAuthorization checks if a user is authorized to access a model
func (s *AuthService) CheckModelAuthorization(userID string, groupIDs []string, model string) (bool, []string, string) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, nil, "user not found"
	}

	if user.Status != "active" {
		return false, nil, "user disabled"
	}

	if len(groupIDs) == 0 {
		return false, nil, "user has no group memberships"
	}

	var allowedModels []string
	for _, groupID := range groupIDs {
		group, err := s.groupRepo.GetByID(groupID)
		if err != nil {
			continue
		}
		if group.TierID == "" {
			continue
		}
		tier, err := s.tierRepo.GetByID(group.TierID)
		if err != nil {
			continue
		}
		allowedModels = append(allowedModels, tier.AllowedModels...)
	}

	if len(allowedModels) == 0 {
		return false, nil, "no tier permissions found"
	}

	for _, pattern := range allowedModels {
		if matchGlob(pattern, model) {
			return true, allowedModels, ""
		}
	}

	return false, allowedModels, "model not in allowed patterns"
}

// Login performs email/password authentication
func (s *AuthService) Login(emailOrUsername, password string) (*entity.User, string, error) {
	var user *entity.User
	var err error

	// Try email first, then username
	log.Printf("[DEBUG] Login attempt with: %s", emailOrUsername)
	user, err = s.userRepo.GetByEmail(emailOrUsername)
	if err != nil {
		log.Printf("[DEBUG] GetByEmail failed, trying username: %v", err)
		user, err = s.userRepo.GetByUsername(emailOrUsername)
	}
	if err != nil {
		log.Printf("[DEBUG] All login attempts failed: %v", err)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if user.Status != "active" {
		return nil, "", fmt.Errorf("user account is disabled")
	}

	if !CheckPassword(user.PasswordHash, password) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := GenerateJWT(user.ID, user.Name, user.Email, user.Role, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// Register creates a new user and returns it
func (s *AuthService) Register(name, email, role, passwordHash string) (*entity.User, error) {
	return s.CreateUser(name, "", email, role, passwordHash)
}

// RegisterWithUsername creates a new user with username
func (s *AuthService) RegisterWithUsername(name, username, email, role, passwordHash string) (*entity.User, error) {
	return s.CreateUser(name, username, email, role, passwordHash)
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(name, username, email, role, passwordHash string) (*entity.User, error) {
	log.Printf("[DEBUG] CreateUser: name=%s, username=%s, email=%s, role=%s", name, username, email, role)

	// Check if user already exists by email
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	user := &entity.User{
		ID:           generateID(),
		Name:         name,
		Username:     username,
		Email:        email,
		Role:         role,
		Status:       "active",
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
// UpdatePassword updates a user's password
func (s *AuthService) UpdatePassword(userID, newPasswordHash string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	user.PasswordHash = newPasswordHash
	return s.userRepo.Update(user)
}

// ResetPassword generates a new random password for a user
func (s *AuthService) ResetPassword(userID string) (string, error) {
	newPass := generateID()[:16]
	hash, err := HashPassword(newPass)
	if err != nil {
		return "", err
	}
	if err := s.UpdatePassword(userID, hash); err != nil {
		return "", err
	}
	return newPass, nil
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
