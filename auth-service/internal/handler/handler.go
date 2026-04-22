package handler

import (
	"context"
	"fmt"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

// Handler implements the gRPC AuthService interface
type Handler struct {
	authv1.UnimplementedAuthServiceServer
	authService *application.AuthService
	userRepo    UserRepository
	apiKeyRepo  APIKeyRepository
}

// UserRepository interface for handler
type UserRepository interface {
	GetByID(id string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.User, int, error)
}

// APIKeyRepository interface for handler
type APIKeyRepository interface {
	GetByID(id string) (*entity.APIKey, error)
	GetByKeyHash(keyHash string) (*entity.APIKey, error)
	GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error)
	Create(apiKey *entity.APIKey) error
	Delete(id string) error
}

// NewHandler creates a new Handler
func NewHandler(authService *application.AuthService, userRepo UserRepository, apiKeyRepo APIKeyRepository) *Handler {
	return &Handler{
		authService: authService,
		userRepo:    userRepo,
		apiKeyRepo:  apiKeyRepo,
	}
}

// ValidateAPIKey validates an API key and returns the user identity
func (h *Handler) ValidateAPIKey(ctx context.Context, req *authv1.ValidateAPIKeyRequest) (*authv1.UserIdentity, error) {
	user, scopes, err := h.authService.ValidateAPIKey(req.ApiKey)
	if err != nil {
		return nil, err
	}

	return &authv1.UserIdentity{
		UserId:   user.ID,
		Role:     user.Role,
		GroupIds: []string{},
		Scopes:   scopes,
	}, nil
}

// Login handles email/password authentication
func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		Token: token,
		User: &authv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Unix(),
		},
	}, nil
}

// CheckModelAuthorization checks if a user is authorized to access a model
func (h *Handler) CheckModelAuthorization(ctx context.Context, req *authv1.CheckModelAuthorizationRequest) (*authv1.AuthorizationResult, error) {
	allowed, models, reason := h.authService.CheckModelAuthorization(req.UserId, req.GroupIds, req.Model)

	return &authv1.AuthorizationResult{
		Allowed:           allowed,
		Reason:            reason,
		AuthorizedModels:  models,
	}, nil
}

// GetUser retrieves a user by ID
func (h *Handler) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.User, error) {
	user, err := h.userRepo.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &authv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

// CreateUser creates a new user
func (h *Handler) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.User, error) {
	var hash string
	if req.Password != "" {
		var err error
		hash, err = application.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
	} else {
		hash, _ = application.HashPassword(req.Password)
	}

	user, err := h.authService.CreateUser(req.Name, req.Email, req.Role, hash)
	if err != nil {
		return nil, err
	}

	return &authv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.User, error) {
	user, err := h.userRepo.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &authv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*commonv1.Empty, error) {
	if err := h.userRepo.Delete(req.Id); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

// ListUsers lists users with pagination
func (h *Handler) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	users, total, err := h.userRepo.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	userProtos := make([]*authv1.User, len(users))
	for i, user := range users {
		userProtos[i] = &authv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Unix(),
		}
	}

	return &authv1.ListUsersResponse{
		Users: userProtos,
		Total: int32(total),
	}, nil
}

// CreateAPIKey creates a new API key for a user
func (h *Handler) CreateAPIKey(ctx context.Context, req *authv1.CreateAPIKeyRequest) (*authv1.CreateAPIKeyResponse, error) {
	keyID, apiKey, err := h.authService.CreateAPIKey(req.UserId, req.Name)
	if err != nil {
		return nil, err
	}

	return &authv1.CreateAPIKeyResponse{
		ApiKeyId: keyID,
		ApiKey:   apiKey,
	}, nil
}

// DeleteAPIKey deletes an API key
func (h *Handler) DeleteAPIKey(ctx context.Context, req *authv1.DeleteAPIKeyRequest) (*commonv1.Empty, error) {
	if err := h.apiKeyRepo.Delete(req.Id); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

// ListAPIKeys lists API keys for a user with pagination
func (h *Handler) ListAPIKeys(ctx context.Context, req *authv1.ListAPIKeysRequest) (*authv1.ListAPIKeysResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	apiKeys, total, err := h.apiKeyRepo.GetByUserID(req.UserId, page, pageSize)
	if err != nil {
		return nil, err
	}

	apiKeyProtos := make([]*authv1.APIKey, len(apiKeys))
	for i, key := range apiKeys {
		var expiresAt int64
		if key.ExpiresAt != nil {
			expiresAt = key.ExpiresAt.Unix()
		}
		apiKeyProtos[i] = &authv1.APIKey{
			Id:        key.ID,
			UserId:    key.UserID,
			Name:      key.Name,
			Scopes:    key.Scopes,
			CreatedAt: key.CreatedAt.Unix(),
			ExpiresAt: expiresAt,
		}
	}

	return &authv1.ListAPIKeysResponse{
		ApiKeys: apiKeyProtos,
		Total:   int32(total),
	}, nil
}

// Phase 2+ Group Management - not implemented for MVP
func (h *Handler) CreateGroup(ctx context.Context, req *authv1.CreateGroupRequest) (*authv1.Group, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) UpdateGroup(ctx context.Context, req *authv1.UpdateGroupRequest) (*authv1.Group, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) DeleteGroup(ctx context.Context, req *authv1.DeleteGroupRequest) (*commonv1.Empty, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) ListGroups(ctx context.Context, req *authv1.ListGroupsRequest) (*authv1.ListGroupsResponse, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) AddUserToGroup(ctx context.Context, req *authv1.AddUserToGroupRequest) (*commonv1.Empty, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) RemoveUserFromGroup(ctx context.Context, req *authv1.RemoveUserFromGroupRequest) (*commonv1.Empty, error) {
	return nil, nil // TODO: implement in Phase 2
}

// Phase 2+ Permission Management - not implemented for MVP
func (h *Handler) GrantPermission(ctx context.Context, req *authv1.GrantPermissionRequest) (*authv1.Permission, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) RevokePermission(ctx context.Context, req *authv1.RevokePermissionRequest) (*commonv1.Empty, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) ListPermissions(ctx context.Context, req *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error) {
	return nil, nil // TODO: implement in Phase 2
}

func (h *Handler) CheckPermission(ctx context.Context, req *authv1.CheckPermissionRequest) (*authv1.CheckPermissionResponse, error) {
	return nil, nil // TODO: implement in Phase 2
}

// Shutdown handles graceful shutdown
func (h *Handler) Shutdown(ctx context.Context) error {
	// Cleanup logic
	return nil
}
