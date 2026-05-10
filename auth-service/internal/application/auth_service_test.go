package application

import (
	"testing"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	users map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*entity.User)}
}

func (m *mockUserRepo) GetByID(id string) (*entity.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) GetByUsername(username string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Create(user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) List(page, pageSize int) ([]*entity.User, int, error) {
	var all []*entity.User
	for _, u := range m.users {
		all = append(all, u)
	}
	total := len(all)
	offset := (page - 1) * pageSize
	if offset >= total {
		return nil, total, nil
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

type mockAPIKeyRepo struct{}

func (m *mockAPIKeyRepo) GetByID(id string) (*entity.APIKey, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAPIKeyRepo) GetByKeyHash(keyHash string) (*entity.APIKey, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAPIKeyRepo) GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error) {
	return nil, 0, nil
}

func (m *mockAPIKeyRepo) Create(apiKey *entity.APIKey) error {
	return nil
}

func (m *mockAPIKeyRepo) Delete(id string) error {
	return nil
}

type mockUserGroupRepoForAuth struct {
	memberships []entity.UserGroupMembership
}

func newMockUserGroupRepoForAuth() *mockUserGroupRepoForAuth {
	return &mockUserGroupRepoForAuth{memberships: make([]entity.UserGroupMembership, 0)}
}

func (m *mockUserGroupRepoForAuth) GetByUserID(userID string) ([]*entity.UserGroupMembership, error) {
	var result []*entity.UserGroupMembership
	for i := range m.memberships {
		if m.memberships[i].UserID == userID {
			result = append(result, &m.memberships[i])
		}
	}
	return result, nil
}

func (m *mockUserGroupRepoForAuth) GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error) {
	var result []*entity.UserGroupMembership
	for i := range m.memberships {
		if m.memberships[i].GroupID == groupID {
			result = append(result, &m.memberships[i])
		}
	}
	return result, len(result), nil
}

func (m *mockUserGroupRepoForAuth) GetGroupIDsByUserID(userID string) ([]string, error) {
	var ids []string
	for _, mem := range m.memberships {
		if mem.UserID == userID {
			ids = append(ids, mem.GroupID)
		}
	}
	return ids, nil
}

func (m *mockUserGroupRepoForAuth) Exists(userID, groupID string) (bool, error) {
	for _, mem := range m.memberships {
		if mem.UserID == userID && mem.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserGroupRepoForAuth) Create(membership *entity.UserGroupMembership) error {
	m.memberships = append(m.memberships, *membership)
	return nil
}

func (m *mockUserGroupRepoForAuth) Delete(userID, groupID string) error {
	for i, mem := range m.memberships {
		if mem.UserID == userID && mem.GroupID == groupID {
			m.memberships = append(m.memberships[:i], m.memberships[i+1:]...)
			return nil
		}
	}
return gorm.ErrRecordNotFound
}

type mockGroupRepoForAuth struct {
	groups map[string]*entity.Group
}

func newMockGroupRepoForAuth() *mockGroupRepoForAuth {
	return &mockGroupRepoForAuth{groups: make(map[string]*entity.Group)}
}

func (m *mockGroupRepoForAuth) Create(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepoForAuth) GetByID(id string) (*entity.Group, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return g, nil
}

func (m *mockGroupRepoForAuth) Update(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepoForAuth) Delete(id string) error {
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepoForAuth) List(page, pageSize int) ([]*entity.Group, int, error) {
	var all []*entity.Group
	for _, g := range m.groups {
		all = append(all, g)
	}
	total := len(all)
	offset := (page - 1) * pageSize
	if offset >= total {
		return nil, total, nil
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

type mockTierRepo struct {
	tiers map[string]*entity.Tier
}

func newMockTierRepo() *mockTierRepo {
	return &mockTierRepo{tiers: make(map[string]*entity.Tier)}
}

func (m *mockTierRepo) GetByID(id string) (*entity.Tier, error) {
	t, ok := m.tiers[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return t, nil
}

func (m *mockTierRepo) Create(tier *entity.Tier) error {
	m.tiers[tier.ID] = tier
	return nil
}

func (m *mockTierRepo) Update(tier *entity.Tier) error {
	m.tiers[tier.ID] = tier
	return nil
}

func (m *mockTierRepo) Delete(id string) error {
	delete(m.tiers, id)
	return nil
}

func (m *mockTierRepo) List(page, pageSize int) ([]*entity.Tier, int, error) {
	var all []*entity.Tier
	for _, t := range m.tiers {
		all = append(all, t)
	}
	return all, len(all), nil
}

func (m *mockTierRepo) GetByName(name string) (*entity.Tier, error) {
	for _, t := range m.tiers {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockTierRepo) GetDefaultTiers() ([]*entity.Tier, error) {
	var result []*entity.Tier
	for _, t := range m.tiers {
		if t.IsDefault {
			result = append(result, t)
		}
	}
	return result, nil
}

func TestAuthService_GenerateAPIKey(t *testing.T) {
	service := &AuthService{}

	// Generate multiple keys to ensure they are unique
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		key, err := service.GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey() error = %v", err)
		}
		if key == "" {
			t.Error("GenerateAPIKey() returned empty string")
		}
		if len(key) < 3 {
			t.Errorf("GenerateAPIKey() returned key shorter than expected: %d", len(key))
		}
		if key[:3] != "sk-" {
			t.Errorf("GenerateAPIKey() returned key with wrong prefix: %s", key[:3])
		}
		if keys[key] {
			t.Error("GenerateAPIKey() generated duplicate key")
		}
		keys[key] = true
	}
}

func TestAuthService_HashAPIKey(t *testing.T) {
	service := &AuthService{}

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "sk-test123456789",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: false, // Hashing empty string is valid
		},
		{
			name:    "long API key",
			apiKey:  "sk-verylongkeywithlotsofrandomcharacters1234567890",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := service.HashAPIKey(tt.apiKey)
			if hash == "" {
				t.Error("HashAPIKey() returned empty string")
			}
			if len(hash) != 64 { // SHA-256 produces 64 hex characters
				t.Errorf("HashAPIKey() returned hash of wrong length: %d", len(hash))
			}

			// Verify that the same key produces the same hash
			hash2 := service.HashAPIKey(tt.apiKey)
			if hash != hash2 {
				t.Error("HashAPIKey() is not deterministic")
			}

			// Verify that different keys produce different hashes
			if tt.apiKey != "" {
				differentKey := tt.apiKey + "different"
				hash3 := service.HashAPIKey(differentKey)
				if hash == hash3 {
					t.Error("HashAPIKey() produced same hash for different keys")
				}
			}
		})
	}
}

func TestAuthService_HashAPIKey_Consistency(t *testing.T) {
	service := &AuthService{}

	apiKey := "sk-test123456789"
	hash1 := service.HashAPIKey(apiKey)
	hash2 := service.HashAPIKey(apiKey)
	hash3 := service.HashAPIKey(apiKey)

	if hash1 != hash2 || hash2 != hash3 {
		t.Error("HashAPIKey() should produce consistent hashes for the same input")
	}
}

func TestAuthService_Login_WithUsername(t *testing.T) {
	userRepo := newMockUserRepo()
	apiKeyRepo := &mockAPIKeyRepo{}
	userGroupRepo := &mockUserGroupRepoForAuth{}
	tierRepo := newMockTierRepo()
	groupRepo := newMockGroupRepoForAuth()

	authSvc := NewAuthService(userRepo, apiKeyRepo, userGroupRepo, tierRepo, groupRepo)

	hash, _ := HashPassword("password123")
	user := &entity.User{
		ID:           "user-1",
		Name:         "Test User",
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         "user",
		Status:       "active",
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	userRepo.users[user.ID] = user

	user, token, err := authSvc.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("Login() with username error = %v", err)
	}
	if token == "" {
		t.Error("Login() with username should return a token")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
}

func TestAuthService_Login_WithEmail(t *testing.T) {
	userRepo := newMockUserRepo()
	apiKeyRepo := &mockAPIKeyRepo{}
	userGroupRepo := &mockUserGroupRepoForAuth{}
	tierRepo := newMockTierRepo()
	groupRepo := newMockGroupRepoForAuth()

	authSvc := NewAuthService(userRepo, apiKeyRepo, userGroupRepo, tierRepo, groupRepo)

	hash, _ := HashPassword("password123")
	user := &entity.User{
		ID:           "user-1",
		Name:         "Test User",
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         "user",
		Status:       "active",
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	userRepo.users[user.ID] = user

	user, token, err := authSvc.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() with email error = %v", err)
	}
	if token == "" {
		t.Error("Login() with email should return a token")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := newMockUserRepo()
	apiKeyRepo := &mockAPIKeyRepo{}
	userGroupRepo := &mockUserGroupRepoForAuth{}
	tierRepo := newMockTierRepo()
	groupRepo := newMockGroupRepoForAuth()

	authSvc := NewAuthService(userRepo, apiKeyRepo, userGroupRepo, tierRepo, groupRepo)

	hash, _ := HashPassword("password123")
	user := &entity.User{
		ID:           "user-1",
		Name:         "Test User",
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         "user",
		Status:       "active",
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	userRepo.users[user.ID] = user

	_, _, err := authSvc.Login("testuser", "wrongpassword")
	if err == nil {
		t.Error("Login() should return error for invalid password")
	}
}

func TestAuthService_CreateUser_WithUsername(t *testing.T) {
	userRepo := newMockUserRepo()
	apiKeyRepo := &mockAPIKeyRepo{}
	userGroupRepo := &mockUserGroupRepoForAuth{}
	tierRepo := newMockTierRepo()
	groupRepo := newMockGroupRepoForAuth()

	authSvc := NewAuthService(userRepo, apiKeyRepo, userGroupRepo, tierRepo, groupRepo)

	hash, _ := HashPassword("password123")
	user, err := authSvc.CreateUser("New User", "newuser", "new@example.com", "user", hash)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if user.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got %s", user.Username)
	}
	if user.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got %s", user.Email)
	}
}

func TestAuthService_CreateUser_DuplicateUsername(t *testing.T) {
	userRepo := newMockUserRepo()
	apiKeyRepo := &mockAPIKeyRepo{}
	userGroupRepo := &mockUserGroupRepoForAuth{}
	tierRepo := newMockTierRepo()
	groupRepo := newMockGroupRepoForAuth()

	authSvc := NewAuthService(userRepo, apiKeyRepo, userGroupRepo, tierRepo, groupRepo)

	hash, _ := HashPassword("password123")
	existingUser := &entity.User{
		ID:           "user-1",
		Name:         "Existing User",
		Username:     "existinguser",
		Email:        "existing@example.com",
		Role:         "user",
		Status:       "active",
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	userRepo.users[existingUser.ID] = existingUser

	_, err := authSvc.CreateUser("New User", "existinguser", "new@example.com", "user", hash)
	if err == nil {
		t.Error("CreateUser() should return error for duplicate username")
	}
}
