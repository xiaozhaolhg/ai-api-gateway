package handler

import "context"

type UserService interface {
	ListUsers(ctx context.Context, page, pageSize int) (*ListUsersResp, error)
	CreateUser(ctx context.Context, name, email, role string) (*User, error)
	UpdateUser(ctx context.Context, id, name, email, role, status string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ListUsersResp struct {
	Users []User `json:"users"`
}

type AdminUsersHandler struct {
	svc UserService
}

func NewAdminUsersHandler(svc UserService) *AdminUsersHandler {
	return &AdminUsersHandler{svc: svc}
}

func (h *AdminUsersHandler) ListUsers(page, pageSize int) (*ListUsersResp, error) {
	return h.svc.ListUsers(context.Background(), page, pageSize)
}

func (h *AdminUsersHandler) CreateUser(name, email, role string) (*User, error) {
	return h.svc.CreateUser(context.Background(), name, email, role)
}

func (h *AdminUsersHandler) UpdateUser(id, name, email, role, status string) (*User, error) {
	return h.svc.UpdateUser(context.Background(), id, name, email, role, status)
}

func (h *AdminUsersHandler) DeleteUser(id string) error {
	return h.svc.DeleteUser(context.Background(), id)
}

type APIKeyService interface {
	ListAPIKeys(ctx context.Context, userID string, page, pageSize int) (*ListAPIKeysResp, error)
	CreateAPIKey(ctx context.Context, userID, name string) (*APIKey, error)
	DeleteAPIKey(ctx context.Context, id string) error
}

type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ListAPIKeysResp struct {
	APIKeys []APIKey `json:"api_keys"`
}

func (h *AdminUsersHandler) ListAPIKeys(userID string, page, pageSize int) (*ListAPIKeysResp, error) {
	return nil, nil
}

func (h *AdminUsersHandler) CreateAPIKey(userID, name string) (*APIKey, error) {
	return nil, nil
}

func (h *AdminUsersHandler) DeleteAPIKey(id string) error {
	return nil
}