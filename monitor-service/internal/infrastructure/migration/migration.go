package migration

import (
	"database/sql"
	"fmt"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createMetricsTable(),
		createAlertRulesTable(),
		createAlertsTable(),
		createProviderHealthTable(),
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}

func createMetricsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS metrics (
		id TEXT PRIMARY KEY,
		provider_id TEXT NOT NULL,
		model TEXT NOT NULL,
		metric_type TEXT NOT NULL,
		value REAL NOT NULL,
		timestamp DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_metrics_provider_id ON metrics(provider_id);
	CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp);
	CREATE INDEX IF NOT EXISTS idx_metrics_type ON metrics(metric_type);
	`
}

func createAlertRulesTable() string {
	return `
	CREATE TABLE IF NOT EXISTS alert_rules (
		id TEXT PRIMARY KEY,
		provider_id TEXT NOT NULL,
		metric_type TEXT NOT NULL,
		threshold REAL NOT NULL,
		operator TEXT NOT NULL,
		window_minutes INTEGER NOT NULL,
		severity TEXT NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_alert_rules_provider_id ON alert_rules(provider_id);
	CREATE INDEX IF NOT EXISTS idx_alert_rules_enabled ON alert_rules(enabled);
	`
}

func createAlertsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS alerts (
		id TEXT PRIMARY KEY,
		alert_rule_id TEXT NOT NULL,
		provider_id TEXT NOT NULL,
		severity TEXT NOT NULL,
		message TEXT NOT NULL,
		value REAL NOT NULL,
		threshold REAL NOT NULL,
		timestamp DATETIME NOT NULL,
		acknowledged BOOLEAN NOT NULL DEFAULT 0,
		acknowledged_at DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_alerts_provider_id ON alerts(provider_id);
	CREATE INDEX IF NOT EXISTS idx_alerts_timestamp ON alerts(timestamp);
	CREATE INDEX IF NOT EXISTS idx_alerts_acknowledged ON alerts(acknowledged);
	`
}

func createProviderHealthTable() string {
	return `
	CREATE TABLE IF NOT EXISTS provider_health (
		provider_id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		latency_ms INTEGER NOT NULL,
		error_rate REAL NOT NULL DEFAULT 0.0,
		last_check_time DATETIME NOT NULL,
		uptime_seconds INTEGER NOT NULL DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_provider_health_status ON provider_health(status);
	CREATE INDEX IF NOT EXISTS idx_provider_health_last_check ON provider_health(last_check_time);
	`
}
