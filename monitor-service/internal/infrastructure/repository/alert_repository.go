package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
	"github.com/google/uuid"
)

// AlertRepository implements the alert repository interface
type AlertRepository struct {
	db *sql.DB
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

// Create creates a new alert
func (r *AlertRepository) Create(alert *entity.Alert) error {
	alert.ID = uuid.New().String()

	query := `
		INSERT INTO alerts (id, alert_rule_id, provider_id, severity, message, value, threshold, timestamp, acknowledged, acknowledged_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, alert.ID, alert.AlertRuleID, alert.ProviderID,
		alert.Severity, alert.Message, alert.Value, alert.Threshold,
		alert.Timestamp, alert.Acknowledged, alert.AcknowledgedAt)
	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	return nil
}

// GetByID retrieves an alert by ID
func (r *AlertRepository) GetByID(id string) (*entity.Alert, error) {
	query := `
		SELECT id, alert_rule_id, provider_id, severity, message, value, threshold, timestamp, acknowledged, acknowledged_at
		FROM alerts
		WHERE id = ?
	`

	var alert entity.Alert
	err := r.db.QueryRow(query, id).Scan(
		&alert.ID, &alert.AlertRuleID, &alert.ProviderID,
		&alert.Severity, &alert.Message, &alert.Value,
		&alert.Threshold, &alert.Timestamp, &alert.Acknowledged, &alert.AcknowledgedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("alert not found")
		}
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	return &alert, nil
}

// GetByProviderID retrieves alerts for a provider
func (r *AlertRepository) GetByProviderID(providerID string, page, pageSize int) ([]*entity.Alert, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE provider_id = ?", providerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	// Get alerts
	query := `
		SELECT id, alert_rule_id, provider_id, severity, message, value, threshold, timestamp, acknowledged, acknowledged_at
		FROM alerts
		WHERE provider_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, providerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query alerts: %w", err)
	}
	defer rows.Close()

	var alerts []*entity.Alert
	for rows.Next() {
		var alert entity.Alert
		err := rows.Scan(
			&alert.ID, &alert.AlertRuleID, &alert.ProviderID,
			&alert.Severity, &alert.Message, &alert.Value,
			&alert.Threshold, &alert.Timestamp, &alert.Acknowledged, &alert.AcknowledgedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, &alert)
	}

	return alerts, total, nil
}

// Acknowledge acknowledges an alert
func (r *AlertRepository) Acknowledge(id string) error {
	now := time.Now()
	query := `UPDATE alerts SET acknowledged = 1, acknowledged_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}
