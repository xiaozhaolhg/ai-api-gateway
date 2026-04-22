package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
	"github.com/google/uuid"
)

// AlertRuleRepository implements the alert rule repository interface
type AlertRuleRepository struct {
	db *sql.DB
}

// NewAlertRuleRepository creates a new alert rule repository
func NewAlertRuleRepository(db *sql.DB) *AlertRuleRepository {
	return &AlertRuleRepository{db: db}
}

// Create creates a new alert rule
func (r *AlertRuleRepository) Create(rule *entity.AlertRule) error {
	rule.ID = uuid.New().String()

	query := `
		INSERT INTO alert_rules (id, provider_id, metric_type, threshold, operator, window_minutes, severity, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, rule.ID, rule.ProviderID, rule.MetricType,
		rule.Threshold, rule.Operator, rule.WindowMinutes, rule.Severity,
		rule.Enabled, rule.CreatedAt, rule.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create alert rule: %w", err)
	}

	return nil
}

// GetByID retrieves an alert rule by ID
func (r *AlertRuleRepository) GetByID(id string) (*entity.AlertRule, error) {
	query := `
		SELECT id, provider_id, metric_type, threshold, operator, window_minutes, severity, enabled, created_at, updated_at
		FROM alert_rules
		WHERE id = ?
	`

	var rule entity.AlertRule
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID, &rule.ProviderID, &rule.MetricType,
		&rule.Threshold, &rule.Operator, &rule.WindowMinutes,
		&rule.Severity, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("alert rule not found")
		}
		return nil, fmt.Errorf("failed to get alert rule: %w", err)
	}

	return &rule, nil
}

// Update updates an existing alert rule
func (r *AlertRuleRepository) Update(rule *entity.AlertRule) error {
	query := `
		UPDATE alert_rules
		SET provider_id = ?, metric_type = ?, threshold = ?, operator = ?, window_minutes = ?, severity = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, rule.ProviderID, rule.MetricType,
		rule.Threshold, rule.Operator, rule.WindowMinutes, rule.Severity,
		rule.Enabled, rule.UpdatedAt, rule.ID)
	if err != nil {
		return fmt.Errorf("failed to update alert rule: %w", err)
	}

	return nil
}

// Delete deletes an alert rule
func (r *AlertRuleRepository) Delete(id string) error {
	query := `DELETE FROM alert_rules WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	return nil
}

// List lists all alert rules
func (r *AlertRuleRepository) List(page, pageSize int) ([]*entity.AlertRule, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM alert_rules").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alert rules: %w", err)
	}

	// Get rules
	query := `
		SELECT id, provider_id, metric_type, threshold, operator, window_minutes, severity, enabled, created_at, updated_at
		FROM alert_rules
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*entity.AlertRule
	for rows.Next() {
		var rule entity.AlertRule
		err := rows.Scan(
			&rule.ID, &rule.ProviderID, &rule.MetricType,
			&rule.Threshold, &rule.Operator, &rule.WindowMinutes,
			&rule.Severity, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	return rules, total, nil
}

// GetEnabledRules retrieves all enabled alert rules
func (r *AlertRuleRepository) GetEnabledRules() ([]*entity.AlertRule, error) {
	query := `
		SELECT id, provider_id, metric_type, threshold, operator, window_minutes, severity, enabled, created_at, updated_at
		FROM alert_rules
		WHERE enabled = 1
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*entity.AlertRule
	for rows.Next() {
		var rule entity.AlertRule
		err := rows.Scan(
			&rule.ID, &rule.ProviderID, &rule.MetricType,
			&rule.Threshold, &rule.Operator, &rule.WindowMinutes,
			&rule.Severity, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}
