package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/migration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting permission to tier migration...")

	dbPath := "/data/auth.db"

	log.Println("Running database migrations...")
	if err := migration.Migrate(dbPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := migratePermissionsToTiers(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func migratePermissionsToTiers(db *gorm.DB) error {
	var permissions []*entity.Permission
	if err := db.Where("resource_type = ? AND effect = ?", "model", "allow").Find(&permissions).Error; err != nil {
		return err
	}

	if len(permissions) == 0 {
		log.Println("No model permissions to migrate")
		return nil
	}

	log.Printf("Found %d model permissions to migrate", len(permissions))

	groupPatterns := make(map[string][]string)
	for _, p := range permissions {
		groupPatterns[p.GroupID] = append(groupPatterns[p.GroupID], p.ResourceID)
	}

	for groupID, patterns := range groupPatterns {
		var existingTier entity.Tier
		if err := db.Where("name = ?", "Migrated-"+groupID).First(&existingTier).Error; err == nil {
			log.Printf("Tier for group %s already exists, skipping", groupID)
			continue
		}

		tier := &entity.Tier{
			ID:               generateID(),
			Name:             "Migrated-" + groupID,
			Description:      "Auto-migrated from permissions",
			IsDefault:        false,
			AllowedModels:    deduplicate(patterns),
			AllowedProviders: extractProviders(patterns),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := db.Create(tier).Error; err != nil {
			log.Printf("Failed to create tier for group %s: %v", groupID, err)
			continue
		}

		if err := db.Model(&entity.Group{}).Where("id = ?", groupID).Update("tier_id", tier.ID).Error; err != nil {
			log.Printf("Failed to assign tier to group %s: %v", groupID, err)
		}

		if err := db.Where("group_id = ? AND resource_type = ? AND effect = ?", groupID, "model", "allow").Delete(&entity.Permission{}).Error; err != nil {
			log.Printf("Failed to delete old permissions for group %s: %v", groupID, err)
		}

		log.Printf("Migrated group %s to tier %s", groupID, tier.ID)
	}

	return nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func deduplicate(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func extractProviders(patterns []string) []string {
	providers := make(map[string]bool)
	for _, pattern := range patterns {
		parts := strings.Split(pattern, ":")
		if len(parts) >= 1 && parts[0] != "" {
			providers[parts[0]] = true
		}
	}
	result := []string{}
	for provider := range providers {
		result = append(result, provider)
	}
	return result
}