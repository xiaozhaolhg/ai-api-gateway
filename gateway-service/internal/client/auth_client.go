package client

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-api-gateway/api/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient wraps the auth-service gRPC client
type AuthClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthClient creates a new auth service client
func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client: authv1.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// Login authenticates a user with email/password and returns a JWT token
func (c *AuthClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	req := &authv1.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return resp, nil
}

// ValidateAPIKey validates an API key and returns user identity
func (c *AuthClient) ValidateAPIKey(ctx context.Context, apiKey string) (*authv1.UserIdentity, error) {
	req := &authv1.ValidateAPIKeyRequest{
		ApiKey: apiKey,
	}

	resp, err := c.client.ValidateAPIKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	return resp, nil
}

// CheckModelAuthorization checks if a user/group is authorized to use a model
func (c *AuthClient) CheckModelAuthorization(ctx context.Context, userID string, groupIDs []string, model string) (*authv1.AuthorizationResult, error) {
	req := &authv1.CheckModelAuthorizationRequest{
		UserId:    userID,
		GroupIds:  groupIDs,
		Model:     model,
	}

	resp, err := c.client.CheckModelAuthorization(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check model authorization: %w", err)
	}

	return resp, nil
}

// GetUser retrieves a user by ID
func (c *AuthClient) GetUser(ctx context.Context, id string) (*authv1.User, error) {
	req := &authv1.GetUserRequest{
		Id: id,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}

// CreateUser creates a new user
func (c *AuthClient) CreateUser(ctx context.Context, name, email, role, password string) (*authv1.User, error) {
	req := &authv1.CreateUserRequest{
		Name:     name,
		Email:    email,
		Role:     role,
		Password: password,
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return resp, nil
}

// UpdateUser updates an existing user
func (c *AuthClient) UpdateUser(ctx context.Context, id, name, email, role, status string) (*authv1.User, error) {
	req := &authv1.UpdateUserRequest{
		Id:     id,
		Name:   name,
		Email:  email,
		Role:   role,
		Status: status,
	}

	resp, err := c.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return resp, nil
}

// DeleteUser deletes a user
func (c *AuthClient) DeleteUser(ctx context.Context, id string) error {
	req := &authv1.DeleteUserRequest{
		Id: id,
	}

	_, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers lists all users
func (c *AuthClient) ListUsers(ctx context.Context, page, pageSize int32) (*authv1.ListUsersResponse, error) {
	req := &authv1.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp, nil
}

// CreateAPIKey creates a new API key
func (c *AuthClient) CreateAPIKey(ctx context.Context, userID, name string) (*authv1.CreateAPIKeyResponse, error) {
	req := &authv1.CreateAPIKeyRequest{
		UserId: userID,
		Name:   name,
	}

	resp, err := c.client.CreateAPIKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return resp, nil
}

// DeleteAPIKey deletes an API key
func (c *AuthClient) DeleteAPIKey(ctx context.Context, id string) error {
	req := &authv1.DeleteAPIKeyRequest{
		Id: id,
	}

	_, err := c.client.DeleteAPIKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

// ListAPIKeys lists API keys for a user
func (c *AuthClient) ListAPIKeys(ctx context.Context, userID string, page, pageSize int32) (*authv1.ListAPIKeysResponse, error) {
	req := &authv1.ListAPIKeysRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListAPIKeys(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	return resp, nil
}
