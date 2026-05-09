package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ai-api-gateway/api/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient wraps the auth-service gRPC client with lazy connection
type AuthClient struct {
	address string
	client  authv1.AuthServiceClient
	conn    *grpc.ClientConn
	mu      sync.RWMutex
}

// NewAuthClient creates a new auth service client with lazy connection
func NewAuthClient(address string) (*AuthClient, error) {
	if address == "" {
		address = "localhost:50051"
	}
	return &AuthClient{
		address: address,
	}, nil
}

// getClient returns the gRPC client, initializing lazily if needed
func (c *AuthClient) getClient() (authv1.AuthServiceClient, error) {
	c.mu.RLock()
	if c.client != nil {
		defer c.mu.RUnlock()
		return c.client, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.client != nil {
		return c.client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(GRPCInterceptor(DefaultRetryConfig())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	c.conn = conn
	c.client = authv1.NewAuthServiceClient(conn)
	return c.client, nil
}

// Close closes the connection
func (c *AuthClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Login authenticates a user with email/password and returns a JWT token
func (c *AuthClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return resp, nil
}

// ValidateAPIKey validates an API key and returns user identity
func (c *AuthClient) ValidateAPIKey(ctx context.Context, apiKey string) (*authv1.UserIdentity, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ValidateAPIKeyRequest{
		ApiKey: apiKey,
	}

	resp, err := client.ValidateAPIKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	return resp, nil
}

// CheckModelAuthorization checks if a user/group is authorized to use a model
func (c *AuthClient) CheckModelAuthorization(ctx context.Context, userID string, groupIDs []string, model string) (*authv1.AuthorizationResult, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.CheckModelAuthorizationRequest{
		UserId:   userID,
		GroupIds: groupIDs,
		Model:    model,
	}

	resp, err := client.CheckModelAuthorization(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check model authorization: %w", err)
	}

	return resp, nil
}

// GetUser retrieves a user by ID
func (c *AuthClient) GetUser(ctx context.Context, id string) (*authv1.User, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.GetUserRequest{
		Id: id,
	}

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}

// CreateUser creates a new user
func (c *AuthClient) CreateUser(ctx context.Context, name, email, role, password string) (*authv1.User, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.CreateUserRequest{
		Name:     name,
		Email:    email,
		Role:     role,
		Password: password,
	}

	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return resp, nil
}

// UpdateUser updates an existing user
func (c *AuthClient) UpdateUser(ctx context.Context, id, name, email, role, status string) (*authv1.User, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.UpdateUserRequest{
		Id:     id,
		Name:   name,
		Email:  email,
		Role:   role,
		Status: status,
	}

	resp, err := client.UpdateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return resp, nil
}

// DeleteUser deletes a user
func (c *AuthClient) DeleteUser(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.DeleteUserRequest{
		Id: id,
	}

	_, err = client.DeleteUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers lists all users
func (c *AuthClient) ListUsers(ctx context.Context, page, pageSize int32) (*authv1.ListUsersResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp, nil
}

// CreateAPIKey creates a new API key
func (c *AuthClient) CreateAPIKey(ctx context.Context, userID, name string) (*authv1.CreateAPIKeyResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.CreateAPIKeyRequest{
		UserId: userID,
		Name:   name,
	}

	resp, err := client.CreateAPIKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return resp, nil
}

// DeleteAPIKey deletes an API key
func (c *AuthClient) DeleteAPIKey(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.DeleteAPIKeyRequest{
		Id: id,
	}

	_, err = client.DeleteAPIKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

// ListAPIKeys lists API keys for a user
func (c *AuthClient) ListAPIKeys(ctx context.Context, userID string, page, pageSize int32) (*authv1.ListAPIKeysResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListAPIKeysRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListAPIKeys(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	return resp, nil
}

// Register creates a new user account
func (c *AuthClient) Register(ctx context.Context, username, email, name, password, role string) (*authv1.RegisterResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.RegisterRequest{
		Username: username,
		Email:    email,
		Name:     name,
		Password: password,
		Role:     role,
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	return resp, nil
}

// Group management methods

// CreateGroup creates a new group
func (c *AuthClient) CreateGroup(ctx context.Context, name, parentGroupID, tierID string) (*authv1.Group, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] CreateGroup gRPC: name=%s, parentGroupID=%s, tierID=%s", name, parentGroupID, tierID)

	req := &authv1.CreateGroupRequest{
		Name:          name,
		ParentGroupId: parentGroupID,
		TierId:        tierID,
	}

	resp, err := client.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	log.Printf("[DEBUG] CreateGroup gRPC response: tierId=%s", resp.TierId)
	return resp, nil
}

// UpdateGroup updates an existing group
func (c *AuthClient) UpdateGroup(ctx context.Context, id, name, parentGroupID, tierID string) (*authv1.Group, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.UpdateGroupRequest{
		Id:            id,
		Name:          name,
		ParentGroupId: parentGroupID,
		TierId:        tierID,
	}

	resp, err := client.UpdateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}
	return resp, nil
}

// DeleteGroup deletes a group
func (c *AuthClient) DeleteGroup(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.DeleteGroupRequest{Id: id}
	_, err = client.DeleteGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

// ListGroups lists groups with pagination
func (c *AuthClient) ListGroups(ctx context.Context, page, pageSize int32) (*authv1.ListGroupsResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListGroupsRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListGroups(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	return resp, nil
}

func (c *AuthClient) ListGroupMembers(ctx context.Context, groupID string, page, pageSize int32) (*authv1.ListGroupMembersResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListGroupMembersRequest{
		GroupId:  groupID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListGroupMembers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list group members: %w", err)
	}
	return resp, nil
}

// AddUserToGroup adds a user to a group
func (c *AuthClient) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.AddUserToGroupRequest{
		UserId:  userID,
		GroupId: groupID,
	}
	_, err = client.AddUserToGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}

// RemoveUserFromGroup removes a user from a group
func (c *AuthClient) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.RemoveUserFromGroupRequest{
		UserId:  userID,
		GroupId: groupID,
	}
	_, err = client.RemoveUserFromGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}
	return nil
}

// GrantPermission grants a permission to a group
func (c *AuthClient) GrantPermission(ctx context.Context, groupID, resourceType, resourceID, action, effect string) (*authv1.Permission, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.GrantPermissionRequest{
		GroupId:      groupID,
		ResourceType: resourceType,
		ResourceId:   resourceID,
		Action:       action,
		Effect:       effect,
	}

	resp, err := client.GrantPermission(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to grant permission: %w", err)
	}

	return resp, nil
}

// RevokePermission revokes a permission
func (c *AuthClient) RevokePermission(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.RevokePermissionRequest{Id: id}
	_, err = client.RevokePermission(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}
	return nil
}

// ListPermissions lists permissions for a group
func (c *AuthClient) ListPermissions(ctx context.Context, groupID string, page, pageSize int32) (*authv1.ListPermissionsResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListPermissionsRequest{
		GroupId:  groupID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListPermissions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return resp, nil
}

// Tier management methods

// CreateTier creates a new tier
func (c *AuthClient) CreateTier(ctx context.Context, name, description string, isDefault bool, allowedModels, allowedProviders []string) (*authv1.Tier, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.CreateTierRequest{
		Name:              name,
		Description:       description,
		IsDefault:         isDefault,
		AllowedModels:     allowedModels,
		AllowedProviders: allowedProviders,
	}

	resp, err := client.CreateTier(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create tier: %w", err)
	}
	return resp, nil
}

// GetTier retrieves a tier by ID
func (c *AuthClient) GetTier(ctx context.Context, id string) (*authv1.Tier, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.GetTierRequest{Id: id}
	resp, err := client.GetTier(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tier: %w", err)
	}
	return resp, nil
}

// UpdateTier updates an existing tier
func (c *AuthClient) UpdateTier(ctx context.Context, id, name, description string, allowedModels, allowedProviders []string) (*authv1.Tier, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.UpdateTierRequest{
		Id:                id,
		Name:              name,
		Description:       description,
		AllowedModels:     allowedModels,
		AllowedProviders:  allowedProviders,
	}

	resp, err := client.UpdateTier(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update tier: %w", err)
	}
	return resp, nil
}

// DeleteTier deletes a tier
func (c *AuthClient) DeleteTier(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.DeleteTierRequest{Id: id}
	_, err = client.DeleteTier(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete tier: %w", err)
	}
	return nil
}

// ListTiers lists tiers with pagination
func (c *AuthClient) ListTiers(ctx context.Context, page, pageSize int32) (*authv1.ListTiersResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &authv1.ListTiersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListTiers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tiers: %w", err)
	}
	return resp, nil
}

// AssignTierToGroup assigns a tier to a group
func (c *AuthClient) AssignTierToGroup(ctx context.Context, groupID, tierID string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.AssignTierToGroupRequest{
		GroupId: groupID,
		TierId:  tierID,
	}
	_, err = client.AssignTierToGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to assign tier to group: %w", err)
	}
	return nil
}

// RemoveTierFromGroup removes a tier from a group
func (c *AuthClient) RemoveTierFromGroup(ctx context.Context, groupID string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &authv1.RemoveTierFromGroupRequest{GroupId: groupID}
	_, err = client.RemoveTierFromGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove tier from group: %w", err)
	}
	return nil
}
