package handler

import (
	"context"
	"testing"
	"time"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
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
	return nil, gorm.ErrRecordNotFound
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

func (m *mockUserRepository) GetByEmail(email string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByUsername(username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
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

// Mock GroupRepository for handler tests
type mockGroupRepositoryForHandler struct {
	groups map[string]*entity.Group
}

func newMockGroupRepositoryForHandler() port.GroupRepository {
	return &mockGroupRepositoryForHandler{groups: make(map[string]*entity.Group)}
}

func (m *mockGroupRepositoryForHandler) Create(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepositoryForHandler) GetByID(id string) (*entity.Group, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return g, nil
}

func (m *mockGroupRepositoryForHandler) Update(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepositoryForHandler) Delete(id string) error {
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepositoryForHandler) List(page, pageSize int) ([]*entity.Group, int, error) {
	var all []*entity.Group
	for _, g := range m.groups {
		all = append(all, g)
	}
	return all, len(all), nil
}

// Mock PermissionRepository for handler tests
type mockPermissionRepositoryForHandler struct {
	permissions map[string]*entity.Permission
}

func newMockPermissionRepositoryForHandler() port.PermissionRepository {
	return &mockPermissionRepositoryForHandler{permissions: make(map[string]*entity.Permission)}
}

func (m *mockPermissionRepositoryForHandler) Create(permission *entity.Permission) error {
	m.permissions[permission.ID] = permission
	return nil
}

func (m *mockPermissionRepositoryForHandler) GetByID(id string) (*entity.Permission, error) {
	p, ok := m.permissions[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (m *mockPermissionRepositoryForHandler) Delete(id string) error {
	delete(m.permissions, id)
	return nil
}

func (m *mockPermissionRepositoryForHandler) ListByGroupID(groupID string, page, pageSize int) ([]*entity.Permission, int, error) {
	var all []*entity.Permission
	for _, p := range m.permissions {
		if p.GroupID == groupID {
			all = append(all, p)
		}
	}
	return all, len(all), nil
}

func (m *mockPermissionRepositoryForHandler) FindByUserGroups(groupIDs []string, resourceType, resourceID, action string) ([]*entity.Permission, error) {
	return nil, nil
}

// Mock UserGroupRepository for handler tests
type mockUserGroupRepositoryForHandler struct {
	memberships []*entity.UserGroupMembership
}

func newMockUserGroupRepositoryForHandler() port.UserGroupRepository {
	return &mockUserGroupRepositoryForHandler{}
}

func (m *mockUserGroupRepositoryForHandler) Create(membership *entity.UserGroupMembership) error {
	m.memberships = append(m.memberships, membership)
	return nil
}

func (m *mockUserGroupRepositoryForHandler) Delete(userID, groupID string) error {
	for i, mg := range m.memberships {
		if mg.UserID == userID && mg.GroupID == groupID {
			m.memberships = append(m.memberships[:i], m.memberships[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockUserGroupRepositoryForHandler) GetByUserID(userID string) ([]*entity.UserGroupMembership, error) {
	var result []*entity.UserGroupMembership
	for _, mg := range m.memberships {
		if mg.UserID == userID {
			result = append(result, mg)
		}
	}
	return result, nil
}

func (m *mockUserGroupRepositoryForHandler) GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error) {
	var result []*entity.UserGroupMembership
	for _, mg := range m.memberships {
		if mg.GroupID == groupID {
			result = append(result, mg)
		}
	}
	return result, len(result), nil
}

func (m *mockUserGroupRepositoryForHandler) Exists(userID, groupID string) (bool, error) {
	for _, mg := range m.memberships {
		if mg.UserID == userID && mg.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}

func setupTestHandler(t *testing.T) *Handler {
	userRepo := newMockUserRepository()
	apiKeyRepo := newMockAPIKeyRepository()
	userGroupRepo := newMockUserGroupRepositoryForHandler()
	authService := application.NewAuthService(userRepo, apiKeyRepo, userGroupRepo)
	groupService := application.NewGroupService(newMockGroupRepositoryForHandler())
	permissionService := application.NewPermissionService(newMockPermissionRepositoryForHandler(), userGroupRepo)
	userGroupService := application.NewUserGroupService(userGroupRepo)

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

	return NewHandler(authService, groupService, permissionService, userGroupService, userRepo, apiKeyRepo)
}

func TestHandler_ValidateAPIKey(t *testing.T) {
	handler := setupTestHandler(t)

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
	handler.apiKeyRepo.(*mockAPIKeyRepository).Create(testAPIKeyRecord)

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

func TestHandler_GroupManagement(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := context.Background()

	// CreateGroup
	resp, err := handler.CreateGroup(ctx, &authv1.CreateGroupRequest{Name: "developers"})
	if err != nil {
		t.Errorf("CreateGroup() error = %v", err)
	}
	if resp.Name != "developers" {
		t.Errorf("Expected name 'developers', got %s", resp.Name)
	}

	// ListGroups
	listResp, err := handler.ListGroups(ctx, &authv1.ListGroupsRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Errorf("ListGroups() error = %v", err)
	}
	if listResp.Total < 1 {
		t.Error("Expected at least 1 group")
	}

	// UpdateGroup
	updateResp, err := handler.UpdateGroup(ctx, &authv1.UpdateGroupRequest{Id: resp.Id, Name: "senior-devs"})
	if err != nil {
		t.Errorf("UpdateGroup() error = %v", err)
	}
	if updateResp.Name != "senior-devs" {
		t.Errorf("Expected name 'senior-devs', got %s", updateResp.Name)
	}

	// DeleteGroup
	_, err = handler.DeleteGroup(ctx, &authv1.DeleteGroupRequest{Id: resp.Id})
	if err != nil {
		t.Errorf("DeleteGroup() error = %v", err)
	}
}

func TestHandler_PermissionManagement(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := context.Background()

	// Create a group first
	group, _ := handler.CreateGroup(ctx, &authv1.CreateGroupRequest{Name: "test-group"})

	// GrantPermission
	perm, err := handler.GrantPermission(ctx, &authv1.GrantPermissionRequest{
		GroupId:      group.Id,
		ResourceType: "model",
		ResourceId:   "gpt-4",
		Action:       "access",
	})
	if err != nil {
		t.Errorf("GrantPermission() error = %v", err)
	}
	if perm.GroupId != group.Id {
		t.Errorf("Expected group_id %s, got %s", group.Id, perm.GroupId)
	}

	// ListPermissions
	listResp, err := handler.ListPermissions(ctx, &authv1.ListPermissionsRequest{
		GroupId:  group.Id,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Errorf("ListPermissions() error = %v", err)
	}
	if listResp.Total < 1 {
		t.Error("Expected at least 1 permission")
	}

	// RevokePermission
	_, err = handler.RevokePermission(ctx, &authv1.RevokePermissionRequest{Id: perm.Id})
	if err != nil {
		t.Errorf("RevokePermission() error = %v", err)
	}
}
