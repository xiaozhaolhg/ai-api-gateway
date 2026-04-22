package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
	"github.com/google/uuid"
)

// MetricRepository implements the metric repository interface
type MetricRepository struct {
	db *sql.DB
}

// NewMetricRepository creates a new metric repository
func NewMetricRepository(db *sql.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// Create creates a new metric
func (r *MetricRepository) Create(metric *entity.Metric) error {
	metric.ID = uuid.New().String()

	query := `
		INSERT INTO metrics (id, provider_id, model, metric_type, value, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, metric.ID, metric.ProviderID, metric.Model,
		metric.MetricType, metric.Value, metric.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create metric: %w", err)
	}

	return nil
}

// GetByProviderID retrieves metrics for a provider
func (r *MetricRepository) GetByProviderID(providerID string, page, pageSize int) ([]*entity.Metric, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM metrics WHERE provider_id = ?", providerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count metrics: %w", err)
	}

	// Get metrics
	query := `
		SELECT id, provider_id, model, metric_type, value, timestamp
		FROM metrics
		WHERE provider_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, providerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*entity.Metric
	for rows.Next() {
		var metric entity.Metric
		err := rows.Scan(
			&metric.ID, &metric.ProviderID, &metric.Model,
			&metric.MetricType, &metric.Value, &metric.Timestamp,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan metric: %w", err)
		}
		metrics = append(metrics, &metric)
	}

	return metrics, total, nil
}

// GetAggregation retrieves aggregated metrics
func (r *MetricRepository) GetAggregation(providerID, metricType, startDate, endDate string) (*entity.MetricAggregation, error) {
	query := `
		SELECT
			provider_id,
			model,
			metric_type,
			AVG(value) as avg_value,
			MIN(value) as min_value,
			MAX(value) as max_value,
			SUM(value) as sum_value,
			COUNT(*) as count
		FROM metrics
		WHERE provider_id = ? AND metric_type = ? AND timestamp BETWEEN ? AND ?
		GROUP BY provider_id, model, metric_type
		LIMIT 1
	`

	var agg entity.MetricAggregation
	err := r.db.QueryRow(query, providerID, metricType, startDate, endDate).Scan(
		&agg.ProviderID, &agg.Model, &agg.MetricType,
		&agg.AvgValue, &agg.MinValue, &agg.MaxValue, &agg.SumValue, &agg.Count,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty aggregation if no records
			return &entity.MetricAggregation{
				ProviderID: providerID,
				MetricType: metricType,
				StartDate:  startDate,
				EndDate:    endDate,
			}, nil
		}
		return nil, fmt.Errorf("failed to get metric aggregation: %w", err)
	}

	agg.StartDate = startDate
	agg.EndDate = endDate

	return &agg, nil
}
