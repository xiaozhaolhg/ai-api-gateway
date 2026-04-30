package client

import (
	"context"
	"fmt"
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
