package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginEndpoint(t *testing.T) {
	r := setupTestRouter()

	body := `{"email":"admin@example.com","password":"password"}`
	req := httptest.NewRequest("POST", "/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "token") {
		t.Error("Expected response to contain token")
	}
}

func TestLogoutEndpoint(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("POST", "/admin/logout", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProtectedEndpointWithoutAuth(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/admin/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProtectedEndpointWithValidToken(t *testing.T) {
	r := setupTestRouter()

	loginReq := httptest.NewRequest("POST", "/admin/login", strings.NewReader(`{"email":"admin@example.com","password":"password"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	cookies := loginW.Result().Cookies()
	var token string
	for _, c := range cookies {
		if c.Name == "auth_token" {
			token = c.Value
			break
		}
	}

	if token == "" {
		t.Fatal("Failed to get token from login")
	}

	meReq := httptest.NewRequest("GET", "/admin/me", nil)
	meReq.Header.Set("Cookie", "auth_token="+token)
	meW := httptest.NewRecorder()

	r.ServeHTTP(meW, meReq)

	if meW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", meW.Code)
	}
}

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/admin/login", handleLogin)
	r.POST("/admin/logout", handleLogout)

	admin := r.Group("/admin")
	admin.Use(jwtAuthMiddleware())
	admin.GET("/me", handleGetCurrentUser)

	return r
}

var _ = TestProtectedEndpointWithValidToken
var _ = TestProtectedEndpointWithoutAuth