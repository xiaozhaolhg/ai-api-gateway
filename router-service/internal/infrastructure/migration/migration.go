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

	err = db.AutoMigrate(&entity.RoutingRule{})
	if err != nil {
		return err
	}

	result := db.Model(&entity.RoutingRule{}).
		Where("is_system_default = ? OR is_system_default IS NULL", false).
		Update("is_system_default", true)
	if result.Error != nil {
		return result.Error
	}

	db.Exec("CREATE INDEX IF NOT EXISTS idx_routing_rules_user_id ON routing_rules(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_routing_rules_model_pattern ON routing_rules(model_pattern)")

	return nil
}
