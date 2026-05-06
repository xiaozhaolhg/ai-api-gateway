package migration

import (
	"os"
	"testing"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigrate_ColumnsAdded(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-router-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = Migrate(tmpFile.Name())
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	var columns []struct {
		Name string
	}
	db.Raw("PRAGMA table_info(routing_rules)").Scan(&columns)

	columnNames := make(map[string]bool)
	for _, col := range columns {
		columnNames[col.Name] = true
	}

	required := []string{"user_id", "is_system_default", "fallback_provider_ids"}
	for _, col := range required {
		if !columnNames[col] {
			t.Errorf("Column %s not found in routing_rules", col)
		}
	}
}

func TestMigrate_NewRowDefaults(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-router-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = Migrate(tmpFile.Name())
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	newRule := entity.RoutingRule{
		ID:           "test-rule-1",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	db.Create(&newRule)

	var fetched entity.RoutingRule
	db.First(&fetched, "id = ?", "test-rule-1")

	if fetched.IsSystemDefault != false {
		t.Errorf("New row IsSystemDefault = %v, want false", fetched.IsSystemDefault)
	}
}

func TestMigrate_IndexesCreated(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-router-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = Migrate(tmpFile.Name())
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	var indexes []struct {
		Name string
	}
	db.Raw("PRAGMA index_list(routing_rules)").Scan(&indexes)

	indexNames := make(map[string]bool)
	for _, idx := range indexes {
		indexNames[idx.Name] = true
	}

	required := []string{"idx_routing_rules_user_id", "idx_routing_rules_model_pattern"}
	for _, idx := range required {
		if !indexNames[idx] {
			t.Errorf("Index %s not found", idx)
		}
	}
}
