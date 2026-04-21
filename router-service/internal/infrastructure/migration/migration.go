package migration

import (
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
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
	err = db.AutoMigrate(&entity.RoutingRule{})
	if err != nil {
		return err
	}

	return nil
}
