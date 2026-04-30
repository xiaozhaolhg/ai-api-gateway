package migration

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&entity.User{}, &entity.APIKey{}, &entity.Group{}, &entity.Permission{}, &entity.UserGroupMembership{})
	if err != nil {
		return err
	}

	return nil
}
