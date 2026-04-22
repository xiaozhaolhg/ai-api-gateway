package handler

import (
	"context"
	"testing"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupLoginTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestLogin_E2E(t *testing.T) {
	db := setupLoginTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := application.NewAuthService(userRepo, nil)
	handler := NewHandler(authService, userRepo, nil)

	password := "securepassword123"
	hash, _ := application.HashPassword(password)

	user := &entity.User{
		ID:           "usr_test123",
		Name:         "Test User",
		Email:        "test@example.com",
		Role:         "admin",
		Status:       "active",
		PasswordHash: hash,
	}
	_ = userRepo.Create(user)

	ctx := context.Background()
	resp, err := handler.Login(ctx, &authv1.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.Token == "" {
		t.Error("Login returned empty token")
	}
	if resp.User.Id != "usr_test123" {
		t.Errorf("Login user.Id = %s, want usr_test123", resp.User.Id)
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("Login user.Email = %s, want test@example.com", resp.User.Email)
	}
	if resp.User.Role != "admin" {
		t.Errorf("Login user.Role = %s, want admin", resp.User.Role)
	}
}

func TestLogin_E2E_InvalidPassword(t *testing.T) {
	db := setupLoginTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := application.NewAuthService(userRepo, nil)
	handler := NewHandler(authService, userRepo, nil)

	hash, _ := application.HashPassword("correctpassword")
	user := &entity.User{
		ID:           "usr_test456",
		Name:         "Test User 2",
		Email:        "user2@example.com",
		Role:         "user",
		Status:       "active",
		PasswordHash: hash,
	}
	_ = userRepo.Create(user)

	ctx := context.Background()
	_, err := handler.Login(ctx, &authv1.LoginRequest{
		Email:    "user2@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("Login should fail for wrong password")
	}
}

func TestLogin_E2E_UserNotFound(t *testing.T) {
	db := setupLoginTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := application.NewAuthService(userRepo, nil)
	handler := NewHandler(authService, userRepo, nil)

	ctx := context.Background()
	_, err := handler.Login(ctx, &authv1.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "anypassword",
	})
	if err == nil {
		t.Fatal("Login should fail for nonexistent user")
	}
}

func TestLogin_E2E_DisabledUser(t *testing.T) {
	db := setupLoginTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := application.NewAuthService(userRepo, nil)
	handler := NewHandler(authService, userRepo, nil)

	hash, _ := application.HashPassword("password123")
	user := &entity.User{
		ID:           "usr_disabled",
		Name:         "Disabled User",
		Email:        "disabled@example.com",
		Role:         "viewer",
		Status:       "disabled",
		PasswordHash: hash,
	}
	_ = userRepo.Create(user)

	ctx := context.Background()
	_, err := handler.Login(ctx, &authv1.LoginRequest{
		Email:    "disabled@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("Login should fail for disabled user")
	}
}

func TestLogin_E2E_Roles(t *testing.T) {
	db := setupLoginTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := application.NewAuthService(userRepo, nil)
	handler := NewHandler(authService, userRepo, nil)

	hash, _ := application.HashPassword("testpass")

	for _, role := range []string{"admin", "user", "viewer"} {
		t.Run(role, func(t *testing.T) {
			email := role + "@test.com"
			user := &entity.User{
				ID:           "usr_" + role,
				Name:         role,
				Email:        email,
				Role:         role,
				Status:       "active",
				PasswordHash: hash,
			}
			_ = userRepo.Create(user)

			ctx := context.Background()
			resp, err := handler.Login(ctx, &authv1.LoginRequest{
				Email:    email,
				Password: "testpass",
			})
			if err != nil {
				t.Errorf("Login failed for role %s: %v", role, err)
			}
			if resp.User.Role != role {
				t.Errorf("Login role = %s, want %s", resp.User.Role, role)
			}
		})
	}
}