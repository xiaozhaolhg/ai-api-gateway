package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
)

// ProviderHealthRepository implements the provider health repository interface
type ProviderHealthRepository struct {
	db *sql.DB
}

// NewProviderHealthRepository creates a new provider health repository
func NewProviderHealthRepository(db *sql.DB) *ProviderHealthRepository {
	return &ProviderHealthRepository{db: db}
}

// Upsert inserts or updates a provider health status
func (r *ProviderHealthRepository) Upsert(status *entity.ProviderHealthStatus) error {
	query := `
		INSERT INTO provider_health (provider_id, status, latency_ms, error_rate, last_check_time, uptime_seconds)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(provider_id) DO UPDATE SET
			status = excluded.status,
			latency_ms = excluded.latency_ms,
			error_rate = excluded.error_rate,
			last_check_time = excluded.last_check_time,
			uptime_seconds = excluded.uptime_seconds
	`

	_, err := r.db.Exec(query, status.ProviderID, status.Status,
		status.LatencyMs, status.ErrorRate, status.LastCheckTime, status.UptimeSeconds)
	if err != nil {
		return fmt.Errorf("failed to upsert provider health: %w", err)
	}

	return nil
}

// GetByProviderID retrieves a provider's health status
func (r *ProviderHealthRepository) GetByProviderID(providerID string) (*entity.ProviderHealthStatus, error) {
	query := `
		SELECT provider_id, status, latency_ms, error_rate, last_check_time, uptime_seconds
		FROM provider_health
		WHERE provider_id = ?
	`

	var status entity.ProviderHealthStatus
	err := r.db.QueryRow(query, providerID).Scan(
		&status.ProviderID, &status.Status, &status.LatencyMs,
		&status.ErrorRate, &status.LastCheckTime, &status.UptimeSeconds,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("provider health not found")
		}
		return nil, fmt.Errorf("failed to get provider health: %w", err)
	}

	return &status, nil
}

// List lists all provider health statuses
func (r *ProviderHealthRepository) List(page, pageSize int) ([]*entity.ProviderHealthStatus, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM provider_health").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count provider health statuses: %w", err)
	}

	// Get statuses
	query := `
		SELECT provider_id, status, latency_ms, error_rate, last_check_time, uptime_seconds
		FROM provider_health
		ORDER BY last_check_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query provider health statuses: %w", err)
	}
	defer rows.Close()

	var statuses []*entity.ProviderHealthStatus
	for rows.Next() {
		var status entity.ProviderHealthStatus
		err := rows.Scan(
			&status.ProviderID, &status.Status, &status.LatencyMs,
			&status.ErrorRate, &status.LastCheckTime, &status.UptimeSeconds,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan provider health status: %w", err)
		}
		statuses = append(statuses, &status)
	}

	return statuses, total, nil
}
