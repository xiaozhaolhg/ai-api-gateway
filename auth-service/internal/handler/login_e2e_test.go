package handler

import (
	"context"
	"testing"

	authv1 "github.com/ai-api-gateway/api/gen/auth/v1"
	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupLoginTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Group{}, &entity.Permission{}, &entity.UserGroupMembership{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func makeTestHandler(db *gorm.DB) (*Handler, port.UserRepository) {
	userRepo := repository.NewUserRepository(db)
	userGroupRepo := repository.NewUserGroupRepository(db)
	tierRepo := repository.NewTierRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	authService := application.NewAuthService(userRepo, nil, userGroupRepo, tierRepo, groupRepo)
	groupService := application.NewGroupService(groupRepo)
	permissionService := application.NewPermissionService(repository.NewPermissionRepository(db), userGroupRepo)
	ugService := application.NewUserGroupService(userGroupRepo)
	tierService := application.NewTierService(tierRepo, groupRepo)
	return NewHandler(authService, groupService, permissionService, ugService, tierService, userRepo, nil, userGroupRepo), userRepo
}

func TestLogin_E2E(t *testing.T) {
	db := setupLoginTestDB(t)
	handler, userRepo := makeTestHandler(db)

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
	handler, userRepo := makeTestHandler(db)

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
	handler, _ := makeTestHandler(db)

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
	handler, userRepo := makeTestHandler(db)

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
	handler, userRepo := makeTestHandler(db)

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

func TestRegister_E2E(t *testing.T) {
	db := setupLoginTestDB(t)
	handler, _ := makeTestHandler(db)

	ctx := context.Background()

	resp, err := handler.Register(ctx, &authv1.RegisterRequest{
		Name:     "New User",
		Email:    "new@example.com",
		Password: "securepassword123",
	})
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
	if resp.User == nil {
		t.Error("Register user = nil, want user")
	}
	if resp.Token == "" {
		t.Error("Register token = empty, want token")
	}
	if resp.User.Name != "New User" {
		t.Errorf("Register name = %s, want New User", resp.User.Name)
	}
}

func TestRegister_E2E_WithUsername(t *testing.T) {
	db := setupLoginTestDB(t)
	handler, _ := makeTestHandler(db)

	ctx := context.Background()

	resp, err := handler.Register(ctx, &authv1.RegisterRequest{
		Name:     "Username User",
		Username: "newuser",
		Password: "securepassword123",
	})
	if err != nil {
		t.Errorf("Register with username failed: %v", err)
	}
	if resp.User == nil {
		t.Error("Register user = nil, want user")
	}
	if resp.User.Email != "newuser@local.dev" {
		t.Errorf("Register email = %s, want newuser@local.dev", resp.User.Email)
	}
}

func TestRegister_E2E_WeakPassword(t *testing.T) {
	db := setupLoginTestDB(t)
	handler, _ := makeTestHandler(db)

	ctx := context.Background()

	_, err := handler.Register(ctx, &authv1.RegisterRequest{
		Name:     "Weak Password User",
		Email:    "weak@example.com",
		Password: "short",
	})
	if err == nil {
		t.Error("Register with weak password should fail")
	}
}

func TestRegister_E2E_Duplicate(t *testing.T) {
	db := setupLoginTestDB(t)
	handler, userRepo := makeTestHandler(db)

	ctx := context.Background()
	hash, _ := application.HashPassword("existingpass")
	existing := &entity.User{
		ID:           "usr_existing",
		Name:         "Existing User",
		Email:        "existing@example.com",
		Role:         "user",
		Status:      "active",
		PasswordHash: hash,
	}
	_ = userRepo.Create(existing)

	_, err := handler.Register(ctx, &authv1.RegisterRequest{
		Name:     "Duplicate User",
		Email:    "existing@example.com",
		Password: "securepassword123",
	})
	if err == nil {
		t.Error("Register with duplicate email should fail")
	}
}