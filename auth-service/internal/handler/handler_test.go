package handler

import (
	"context"
	"testing"
	"time"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

// MockUserRepository for testing
type mockUserRepository struct {
	users map[string]*entity.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *mockUserRepository) GetByID(id string) (*entity.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, nil // Simulate not found
}

func (m *mockUserRepository) Create(user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Update(user *entity.User) error {
	if _, ok := m.users[user.ID]; ok {
		m.users[user.ID] = user
		return nil
	}
	return nil // Simulate not found
}

func (m *mockUserRepository) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) List(page, pageSize int) ([]*entity.User, int, error) {
	users := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

// MockAPIKeyRepository for testing
type mockAPIKeyRepository struct {
	apiKeys map[string]*entity.APIKey
}

func newMockAPIKeyRepository() *mockAPIKeyRepository {
	return &mockAPIKeyRepository{
		apiKeys: make(map[string]*entity.APIKey),
	}
}

func (m *mockAPIKeyRepository) GetByID(id string) (*entity.APIKey, error) {
	if key, ok := m.apiKeys[id]; ok {
		return key, nil
	}
	return nil, nil // Simulate not found
}

func (m *mockAPIKeyRepository) GetByKeyHash(keyHash string) (*entity.APIKey, error) {
	for _, key := range m.apiKeys {
		if key.KeyHash == keyHash {
			return key, nil
		}
	}
	return nil, nil // Simulate not found
}

func (m *mockAPIKeyRepository) GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error) {
	keys := make([]*entity.APIKey, 0)
	for _, key := range m.apiKeys {
		if key.UserID == userID {
			keys = append(keys, key)
		}
	}
	return keys, len(keys), nil
}

func (m *mockAPIKeyRepository) Create(apiKey *entity.APIKey) error {
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}

func (m *mockAPIKeyRepository) Delete(id string) error {
	delete(m.apiKeys, id)
	return nil
}

func setupTestHandler(t *testing.T) *Handler {
	userRepo := newMockUserRepository()
	apiKeyRepo := newMockAPIKeyRepository()
	authService := application.NewAuthService(userRepo, apiKeyRepo)

	// Create a test user
	user := &entity.User{
		ID:        "test-user-1",
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Create a test API key
	apiKey, _ := authService.GenerateAPIKey()
	apiKeyRecord := &entity.APIKey{
		ID:        "test-key-1",
		UserID:    user.ID,
		KeyHash:   authService.HashAPIKey(apiKey),
		Name:      "Test Key",
		Scopes:    []string{"read", "write"},
		CreatedAt: time.Now(),
	}
	apiKeyRepo.Create(apiKeyRecord)

	return NewHandler(authService, userRepo, apiKeyRepo)
}

func TestHandler_ValidateAPIKey(t *testing.T) {
	handler := setupTestHandler(t)

	// Get the test API key
	mockRepo := handler.apiKeyRepo.(*mockAPIKeyRepository)
	testKey := mockRepo.apiKeys["test-key-1"]
	testAPIKey := "sk-test123" // This won't match the hash, so we need to use the actual key

	// For this test, we'll create a new key with a known value
	apiKey, _ := handler.authService.GenerateAPIKey()
	hash := handler.authService.HashAPIKey(apiKey)
	testAPIKeyRecord := &entity.APIKey{
		ID:        "known-key-1",
		UserID:    "test-user-1",
		KeyHash:   hash,
		Name:      "Known Key",
		Scopes:    []string{"read"},
		CreatedAt: time.Now(),
	}
	mockRepo.Create(testAPIKeyRecord)

	req := &authv1.ValidateAPIKeyRequest{
		ApiKey: apiKey,
	}

	resp, err := handler.ValidateAPIKey(context.Background(), req)
	if err != nil {
		t.Errorf("ValidateAPIKey() error = %v", err)
	}
	if resp.UserId != "test-user-1" {
		t.Errorf("Expected user ID test-user-1, got %s", resp.UserId)
	}
	if len(resp.Scopes) != 1 || resp.Scopes[0] != "read" {
		t.Errorf("Expected scopes [read], got %v", resp.Scopes)
	}

	// Test invalid key
	invalidReq := &authv1.ValidateAPIKeyRequest{
		ApiKey: "invalid-key",
	}
	_, err = handler.ValidateAPIKey(context.Background(), invalidReq)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}
}

func TestHandler_CheckModelAuthorization(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.CheckModelAuthorizationRequest{
		UserId:    "test-user-1",
		GroupIds:  []string{},
		Model:     "ollama:llama2",
	}

	resp, err := handler.CheckModelAuthorization(context.Background(), req)
	if err != nil {
		t.Errorf("CheckModelAuthorization() error = %v", err)
	}
	if !resp.Allowed {
		t.Error("Expected authorization to be allowed for MVP")
	}
	if len(resp.AuthorizedModels) == 0 {
		t.Error("Expected authorized models to be non-empty")
	}
}

func TestHandler_GetUser(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.GetUserRequest{
		Id: "test-user-1",
	}

	resp, err := handler.GetUser(context.Background(), req)
	if err != nil {
		t.Errorf("GetUser() error = %v", err)
	}
	if resp.Id != "test-user-1" {
		t.Errorf("Expected user ID test-user-1, got %s", resp.Id)
	}
	if resp.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", resp.Name)
	}
	if resp.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", resp.Email)
	}

	// Test non-existent user
	invalidReq := &authv1.GetUserRequest{
		Id: "non-existent",
	}
	_, err = handler.GetUser(context.Background(), invalidReq)
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

func TestHandler_CreateUser(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.CreateUserRequest{
		Name:  "New User",
		Email: "newuser@example.com",
		Role:  "user",
	}

	resp, err := handler.CreateUser(context.Background(), req)
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
	}
	if resp.Name != "New User" {
		t.Errorf("Expected name 'New User', got %s", resp.Name)
	}
	if resp.Email != "newuser@example.com" {
		t.Errorf("Expected email 'newuser@example.com', got %s", resp.Email)
	}
	if resp.Role != "user" {
		t.Errorf("Expected role 'user', got %s", resp.Role)
	}
	if resp.Status != "active" {
		t.Errorf("Expected status 'active', got %s", resp.Status)
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.UpdateUserRequest{
		Id:     "test-user-1",
		Name:   "Updated Name",
		Status: "disabled",
	}

	resp, err := handler.UpdateUser(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateUser() error = %v", err)
	}
	if resp.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", resp.Name)
	}
	if resp.Status != "disabled" {
		t.Errorf("Expected status 'disabled', got %s", resp.Status)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.DeleteUserRequest{
		Id: "test-user-1",
	}

	_, err := handler.DeleteUser(context.Background(), req)
	if err != nil {
		t.Errorf("DeleteUser() error = %v", err)
	}

	// Verify deletion
	getReq := &authv1.GetUserRequest{
		Id: "test-user-1",
	}
	_, err = handler.GetUser(context.Background(), getReq)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestHandler_ListUsers(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := handler.ListUsers(context.Background(), req)
	if err != nil {
		t.Errorf("ListUsers() error = %v", err)
	}
	if resp.Total < 1 {
		t.Error("Expected at least 1 user")
	}
	if len(resp.Users) < 1 {
		t.Error("Expected at least 1 user in response")
	}
}

func TestHandler_CreateAPIKey(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.CreateAPIKeyRequest{
		UserId: "test-user-1",
		Name:   "New Key",
	}

	resp, err := handler.CreateAPIKey(context.Background(), req)
	if err != nil {
		t.Errorf("CreateAPIKey() error = %v", err)
	}
	if resp.ApiKeyId == "" {
		t.Error("Expected API key ID to be non-empty")
	}
	if resp.ApiKey == "" {
		t.Error("Expected API key to be non-empty")
	}
	if len(resp.ApiKey) < 3 || resp.ApiKey[:3] != "sk-" {
		t.Error("Expected API key to start with 'sk-'")
	}
}

func TestHandler_DeleteAPIKey(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.DeleteAPIKeyRequest{
		Id: "test-key-1",
	}

	_, err := handler.DeleteAPIKey(context.Background(), req)
	if err != nil {
		t.Errorf("DeleteAPIKey() error = %v", err)
	}

	// Verify deletion
	mockRepo := handler.apiKeyRepo.(*mockAPIKeyRepository)
	if _, ok := mockRepo.apiKeys["test-key-1"]; ok {
		t.Error("Expected API key to be deleted")
	}
}

func TestHandler_ListAPIKeys(t *testing.T) {
	handler := setupTestHandler(t)

	req := &authv1.ListAPIKeysRequest{
		UserId:  "test-user-1",
		Page:    1,
		PageSize: 10,
	}

	resp, err := handler.ListAPIKeys(context.Background(), req)
	if err != nil {
		t.Errorf("ListAPIKeys() error = %v", err)
	}
	if resp.Total < 0 {
		t.Error("Expected total to be >= 0")
	}
}

func TestHandler_Phase2Methods(t *testing.T) {
	handler := setupTestHandler(t)

	// Test that Phase 2+ methods return nil (not implemented)
	ctx := context.Background()

	_, err := handler.CreateGroup(ctx, &authv1.CreateGroupRequest{})
	if err != nil {
		t.Errorf("CreateGroup should not error, got %v", err)
	}

	_, err = handler.UpdateGroup(ctx, &authv1.UpdateGroupRequest{})
	if err != nil {
		t.Errorf("UpdateGroup should not error, got %v", err)
	}

	_, err = handler.DeleteGroup(ctx, &authv1.DeleteGroupRequest{})
	if err != nil {
		t.Errorf("DeleteGroup should not error, got %v", err)
	}

	_, err = handler.ListGroups(ctx, &authv1.ListGroupsRequest{})
	if err != nil {
		t.Errorf("ListGroups should not error, got %v", err)
	}

	_, err = handler.GrantPermission(ctx, &authv1.GrantPermissionRequest{})
	if err != nil {
		t.Errorf("GrantPermission should not error, got %v", err)
	}

	_, err = handler.CheckPermission(ctx, &authv1.CheckPermissionRequest{})
	if err != nil {
		t.Errorf("CheckPermission should not error, got %v", err)
	}
}
