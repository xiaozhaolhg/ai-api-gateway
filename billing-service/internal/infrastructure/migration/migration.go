package migration

import (
	"database/sql"
	"fmt"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsageRecordsTable(),
		createPricingRulesTable(),
		createBillingAccountsTable(),
		createBudgetsTable(),
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}

func createUsageRecordsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS usage_records (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		provider_id TEXT NOT NULL,
		model TEXT NOT NULL,
		prompt_tokens INTEGER NOT NULL,
		completion_tokens INTEGER NOT NULL,
		cost REAL NOT NULL,
		timestamp DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_usage_records_user_id ON usage_records(user_id);
	CREATE INDEX IF NOT EXISTS idx_usage_records_timestamp ON usage_records(timestamp);
	`
}

func createPricingRulesTable() string {
	return `
	CREATE TABLE IF NOT EXISTS pricing_rules (
		id TEXT PRIMARY KEY,
		provider_id TEXT NOT NULL,
		model TEXT NOT NULL,
		prompt_price_per_1k REAL NOT NULL,
		completion_price_per_1k REAL NOT NULL,
		currency TEXT NOT NULL DEFAULT 'USD',
		effective_from DATETIME NOT NULL,
		effective_until DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_pricing_rules_provider_model ON pricing_rules(provider_id, model);
	CREATE INDEX IF NOT EXISTS idx_pricing_rules_effective_from ON pricing_rules(effective_from);
	`
}

func createBillingAccountsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS billing_accounts (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL UNIQUE,
		balance REAL NOT NULL DEFAULT 0.0,
		currency TEXT NOT NULL DEFAULT 'USD',
		status TEXT NOT NULL DEFAULT 'active',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_billing_accounts_user_id ON billing_accounts(user_id);
	`
}

func createBudgetsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS budgets (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		limit REAL NOT NULL,
		period TEXT NOT NULL,
		soft_cap REAL NOT NULL,
		hard_cap REAL NOT NULL,
		start_date DATETIME NOT NULL,
		end_date DATETIME,
		alert_threshold REAL NOT NULL DEFAULT 0.8,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);
	CREATE INDEX IF NOT EXISTS idx_budgets_dates ON budgets(start_date, end_date);
	`
}
