package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/application"
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/migration"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

func main() {
	dbPath := "/data/auth.db"

	// Run migrations
	log.Println("Running database migrations...")
	if err := migration.Migrate(dbPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	// Check if admin user already exists
	log.Println("Checking for existing admin user...")
	existingUsers, _, err := userRepo.List(1, 10)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	if len(existingUsers) > 0 {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	// Create default admin user
	log.Println("Creating default admin user...")
	passwordHash, err := application.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}
	adminUser := &entity.User{
		ID:           generateID(),
		Username:     "admin",
		Name:         "Admin",
		Email:        "admin@example.com",
		Role:         "admin",
		Status:       "active",
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	if err := userRepo.Create(adminUser); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}
	log.Printf("Created admin user: %s (ID: %s)", adminUser.Name, adminUser.ID)

	// Create API key for admin user
	log.Println("Creating API key for admin user...")
	apiKey := &entity.APIKey{
		ID:        generateID(),
		UserID:    adminUser.ID,
		KeyHash:   hashAPIKey("sk-admin-default-key-12345"),
		Name:      "Default Admin Key",
		Scopes:    []string{"read", "write", "admin"},
		CreatedAt: time.Now(),
	}

	if err := apiKeyRepo.Create(apiKey); err != nil {
		log.Fatalf("Failed to create API key: %v", err)
	}
	log.Printf("Created API key: %s (Key: sk-admin-default-key-12345)", apiKey.Name)

	log.Println("Seed completed successfully!")
	log.Println("Admin credentials:")
	log.Println("  Username: admin")
	log.Println("  Password: admin123")
	log.Println("  API Key: sk-admin-default-key-12345")
}
