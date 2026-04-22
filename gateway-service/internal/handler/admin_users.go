package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// AdminUsersHandler handles admin user management requests
type AdminUsersHandler struct {
	authClient *client.AuthClient
}

// NewAdminUsersHandler creates a new admin users handler
func NewAdminUsersHandler(authClient *client.AuthClient) *AdminUsersHandler {
	return &AdminUsersHandler{
		authClient: authClient,
	}
}

// ServeHTTP handles HTTP requests for /admin/users and /admin/api-keys
func (h *AdminUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/admin/api-keys") {
		h.handleAPIKeys(w, r)
	} else {
		h.handleUsers(w, r)
	}
}

// handleUsers handles user management
func (h *AdminUsersHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodPut:
		h.updateUser(w, r)
	case http.MethodDelete:
		h.deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleAPIKeys handles API key management
func (h *AdminUsersHandler) handleAPIKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAPIKeys(w, r)
	case http.MethodPost:
		h.createAPIKey(w, r)
	case http.MethodDelete:
		h.deleteAPIKey(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listUsers lists all users
func (h *AdminUsersHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	resp, err := h.authClient.ListUsers(r.Context(), int32(page), int32(pageSize))
	if err != nil {
		http.Error(w, "Failed to list users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createUser creates a new user
func (h *AdminUsersHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.CreateUser(r.Context(), reqBody.Name, reqBody.Email, reqBody.Role)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// updateUser updates an existing user
func (h *AdminUsersHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.UpdateUser(r.Context(), id, reqBody.Name, reqBody.Email, reqBody.Role, reqBody.Status)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// deleteUser deletes a user
func (h *AdminUsersHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}

	err := h.authClient.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listAPIKeys lists API keys for a user
func (h *AdminUsersHandler) listAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	resp, err := h.authClient.ListAPIKeys(r.Context(), userID, int32(page), int32(pageSize))
	if err != nil {
		http.Error(w, "Failed to list API keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createAPIKey creates a new API key
func (h *AdminUsersHandler) createAPIKey(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserId string `json:"user_id"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.CreateAPIKey(r.Context(), reqBody.UserId, reqBody.Name)
	if err != nil {
		http.Error(w, "Failed to create API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// deleteAPIKey deletes an API key
func (h *AdminUsersHandler) deleteAPIKey(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing API key id", http.StatusBadRequest)
		return
	}

	err := h.authClient.DeleteAPIKey(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
