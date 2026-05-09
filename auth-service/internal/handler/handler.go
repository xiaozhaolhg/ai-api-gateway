package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

// Handler implements the gRPC AuthService interface
type Handler struct {
	authv1.UnimplementedAuthServiceServer
	authService       *application.AuthService
	groupService      *application.GroupService
	permissionService *application.PermissionService
	userGroupService  *application.UserGroupService
	tierService       *application.TierService
	userRepo          UserRepository
	apiKeyRepo        APIKeyRepository
}

// UserRepository interface for handler
type UserRepository interface {
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
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
func NewHandler(authService *application.AuthService, groupService *application.GroupService, permissionService *application.PermissionService, userGroupService *application.UserGroupService, tierService *application.TierService, userRepo UserRepository, apiKeyRepo APIKeyRepository) *Handler {
	return &Handler{
		authService:       authService,
		groupService:      groupService,
		permissionService: permissionService,
		userGroupService:  userGroupService,
		tierService:       tierService,
		userRepo:          userRepo,
		apiKeyRepo:        apiKeyRepo,
	}
}

// ValidateAPIKey validates an API key and returns the user identity
func (h *Handler) ValidateAPIKey(ctx context.Context, req *authv1.ValidateAPIKeyRequest) (*authv1.UserIdentity, error) {
	user, scopes, groupIDs, err := h.authService.ValidateAPIKey(req.ApiKey)
	if err != nil {
		return nil, err
	}

	return &authv1.UserIdentity{
		UserId:   user.ID,
		Role:     user.Role,
		GroupIds: groupIDs,
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

// Register handles user registration with username or email
func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	email := req.Email
	if email == "" && req.Username != "" {
		email = req.Username + "@local.dev"
	}

	existing, _ := h.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Check username uniqueness if provided
	if req.Username != "" {
		existingUsername, _ := h.userRepo.GetByUsername(req.Username)
		if existingUsername != nil {
			return nil, fmt.Errorf("username already taken")
		}
	}

	hash, err := application.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	log.Printf("[DEBUG] Register: name=%s, username=%s, email=%s, role=%s", req.Name, req.Username, email, role)

	user, err := h.authService.RegisterWithUsername(req.Name, req.Username, email, role, hash)
	if err != nil {
		return nil, err
	}

	token, err := application.GenerateJWT(user.ID, user.Name, user.Email, user.Role, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &authv1.RegisterResponse{
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

	user, err := h.authService.CreateUser(req.Name, req.Username, req.Email, req.Role, hash)
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

// Group Management

func (h *Handler) CreateGroup(ctx context.Context, req *authv1.CreateGroupRequest) (*authv1.Group, error) {
	log.Printf("[DEBUG] CreateGroup: name=%s, parentGroupId=%s, TierId=%s", req.Name, req.ParentGroupId, req.TierId)

	group, err := h.groupService.CreateGroup(req.Name, "", req.ParentGroupId, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Group created: id=%s", group.ID)

	if req.TierId != "" {
		log.Printf("[DEBUG] Attempting to assign tier %s to group %s", req.TierId, group.ID)
		if err := h.tierService.AssignTierToGroup(group.ID, req.TierId); err != nil {
			log.Printf("[ERROR] Failed to assign tier: %v", err)
		}
		group, _ = h.groupService.GetGroupByID(group.ID)
		log.Printf("[DEBUG] Group after tier assignment: tierId=%s", group.TierID)
	}

	return &authv1.Group{
		Id:            group.ID,
		Name:          group.Name,
		Description:   group.Description,
		ParentGroupId: group.ParentGroupID,
		TierId:        group.TierID,
		CreatedAt:     group.CreatedAt.Unix(),
		UpdatedAt:     group.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) UpdateGroup(ctx context.Context, req *authv1.UpdateGroupRequest) (*authv1.Group, error) {
	group, err := h.groupService.UpdateGroup(req.Id, req.Name, req.ParentGroupId)
	if err != nil {
		return nil, err
	}

	if req.TierId != "" {
		if err := h.tierService.AssignTierToGroup(req.Id, req.TierId); err != nil {
			log.Printf("[ERROR] Failed to assign tier: %v", err)
		}
	} else {
		_ = h.tierService.RemoveTierFromGroup(req.Id)
	}
	group, _ = h.groupService.GetGroupByID(req.Id)

	return &authv1.Group{
		Id:            group.ID,
		Name:          group.Name,
		Description:   group.Description,
		ParentGroupId: group.ParentGroupID,
		TierId:        group.TierID,
		CreatedAt:     group.CreatedAt.Unix(),
		UpdatedAt:     group.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) DeleteGroup(ctx context.Context, req *authv1.DeleteGroupRequest) (*commonv1.Empty, error) {
	if err := h.groupService.DeleteGroup(req.Id); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) ListGroups(ctx context.Context, req *authv1.ListGroupsRequest) (*authv1.ListGroupsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	groups, total, err := h.groupService.ListGroups(page, pageSize)
	if err != nil {
		return nil, err
	}

	groupProtos := make([]*authv1.Group, len(groups))
	for i, group := range groups {
		groupProtos[i] = &authv1.Group{
			Id:            group.ID,
			Name:          group.Name,
			Description:   group.Description,
			ParentGroupId: group.ParentGroupID,
			TierId:        group.TierID,
			CreatedAt:     group.CreatedAt.Unix(),
			UpdatedAt:     group.UpdatedAt.Unix(),
		}
	}

	return &authv1.ListGroupsResponse{
		Groups: groupProtos,
		Total:  int32(total),
	}, nil
}

func (h *Handler) ListGroupMembers(ctx context.Context, req *authv1.ListGroupMembersRequest) (*authv1.ListGroupMembersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	memberships, total, err := h.userGroupService.GetGroupMembers(req.GroupId, page, pageSize)
	if err != nil {
		return nil, err
	}

	memberProtos := make([]*authv1.User, 0, len(memberships))
	for _, m := range memberships {
		user, err := h.userRepo.GetByID(m.UserID)
		if err != nil {
			continue
		}
		memberProtos = append(memberProtos, &authv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Unix(),
		})
	}

	return &authv1.ListGroupMembersResponse{
		Members: memberProtos,
		Total:   int32(total),
	}, nil
}

func (h *Handler) AddUserToGroup(ctx context.Context, req *authv1.AddUserToGroupRequest) (*commonv1.Empty, error) {
	if err := h.userGroupService.AddUserToGroup(req.UserId, req.GroupId); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) RemoveUserFromGroup(ctx context.Context, req *authv1.RemoveUserFromGroupRequest) (*commonv1.Empty, error) {
	if err := h.userGroupService.RemoveUserFromGroup(req.UserId, req.GroupId); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

// Permission Management

func (h *Handler) GrantPermission(ctx context.Context, req *authv1.GrantPermissionRequest) (*authv1.Permission, error) {
	permission, err := h.permissionService.GrantPermission(req.GroupId, req.ResourceType, req.ResourceId, req.Action, "allow")
	if err != nil {
		return nil, err
	}

return &authv1.Permission{
		Id:           permission.ID,
		GroupId:      permission.GroupID,
		ResourceType: permission.ResourceType,
		ResourceId:   permission.ResourceID,
		Action:       permission.Action,
		Effect:       permission.Effect,
		CreatedAt:    permission.CreatedAt.Unix(),
		UpdatedAt:    permission.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) RevokePermission(ctx context.Context, req *authv1.RevokePermissionRequest) (*commonv1.Empty, error) {
	if err := h.permissionService.RevokePermission(req.Id); err != nil {
		return nil, err
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) ListPermissions(ctx context.Context, req *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	permissions, total, err := h.permissionService.ListPermissions(req.GroupId, page, pageSize)
	if err != nil {
		return nil, err
	}

	permissionProtos := make([]*authv1.Permission, len(permissions))
	for i, p := range permissions {
		permissionProtos[i] = &authv1.Permission{
			Id:           p.ID,
			GroupId:      p.GroupID,
			ResourceType: p.ResourceType,
			ResourceId:   p.ResourceID,
			Action:       p.Action,
			Effect:       p.Effect,
			CreatedAt:    p.CreatedAt.Unix(),
			UpdatedAt:    p.UpdatedAt.Unix(),
		}
	}

	return &authv1.ListPermissionsResponse{
		Permissions: permissionProtos,
		Total:       int32(total),
	}, nil
}

func (h *Handler) CheckPermission(ctx context.Context, req *authv1.CheckPermissionRequest) (*authv1.CheckPermissionResponse, error) {
	allowed, err := h.permissionService.CheckPermission(req.UserId, req.ResourceType, req.ResourceId, req.Action)
	if err != nil {
		return nil, err
	}

	return &authv1.CheckPermissionResponse{
		Allowed: allowed,
	}, nil
}

func (h *Handler) CreateTier(ctx context.Context, req *authv1.CreateTierRequest) (*authv1.Tier, error) {
	tier, err := h.tierService.CreateTier(req.Name, req.Description, req.IsDefault, req.AllowedModels, req.AllowedProviders)
	if err != nil {
		return nil, err
	}

	return &authv1.Tier{
		Id:                tier.ID,
		Name:              tier.Name,
		Description:       tier.Description,
		IsDefault:         tier.IsDefault,
		AllowedModels:     tier.AllowedModels,
		AllowedProviders:  tier.AllowedProviders,
		CreatedAt:         tier.CreatedAt.Unix(),
		UpdatedAt:         tier.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) GetTier(ctx context.Context, req *authv1.GetTierRequest) (*authv1.Tier, error) {
	tier, err := h.tierService.GetTier(req.Id)
	if err != nil {
		return nil, err
	}

	return &authv1.Tier{
		Id:                tier.ID,
		Name:              tier.Name,
		Description:       tier.Description,
		IsDefault:         tier.IsDefault,
		AllowedModels:     tier.AllowedModels,
		AllowedProviders:  tier.AllowedProviders,
		CreatedAt:         tier.CreatedAt.Unix(),
		UpdatedAt:         tier.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) UpdateTier(ctx context.Context, req *authv1.UpdateTierRequest) (*authv1.Tier, error) {
	tier, err := h.tierService.UpdateTier(req.Id, req.Name, req.Description, req.AllowedModels, req.AllowedProviders)
	if err != nil {
		return nil, err
	}

	return &authv1.Tier{
		Id:                tier.ID,
		Name:              tier.Name,
		Description:       tier.Description,
		IsDefault:         tier.IsDefault,
		AllowedModels:     tier.AllowedModels,
		AllowedProviders:  tier.AllowedProviders,
		CreatedAt:         tier.CreatedAt.Unix(),
		UpdatedAt:         tier.UpdatedAt.Unix(),
	}, nil
}

func (h *Handler) DeleteTier(ctx context.Context, req *authv1.DeleteTierRequest) (*commonv1.Empty, error) {
	if err := h.tierService.DeleteTier(req.Id); err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

func (h *Handler) ListTiers(ctx context.Context, req *authv1.ListTiersRequest) (*authv1.ListTiersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	tiers, total, err := h.tierService.ListTiers(page, pageSize)
	if err != nil {
		return nil, err
	}

	tierProtos := make([]*authv1.Tier, len(tiers))
	for i, t := range tiers {
		tierProtos[i] = &authv1.Tier{
			Id:                t.ID,
			Name:              t.Name,
			Description:       t.Description,
			IsDefault:         t.IsDefault,
			AllowedModels:     t.AllowedModels,
			AllowedProviders:  t.AllowedProviders,
			CreatedAt:         t.CreatedAt.Unix(),
			UpdatedAt:         t.UpdatedAt.Unix(),
		}
	}

	return &authv1.ListTiersResponse{
		Tiers: tierProtos,
		Total: int32(total),
	}, nil
}

func (h *Handler) AssignTierToGroup(ctx context.Context, req *authv1.AssignTierToGroupRequest) (*commonv1.Empty, error) {
	if err := h.tierService.AssignTierToGroup(req.GroupId, req.TierId); err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

func (h *Handler) RemoveTierFromGroup(ctx context.Context, req *authv1.RemoveTierFromGroupRequest) (*commonv1.Empty, error) {
	if err := h.tierService.RemoveTierFromGroup(req.GroupId); err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	// Cleanup logic
	return nil
}
