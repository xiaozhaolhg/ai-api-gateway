package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UsageRecordRepository implements the usage record repository interface
type UsageRecordRepository struct {
	db *sql.DB
}

// NewUsageRecordRepository creates a new usage record repository
func NewUsageRecordRepository(db *sql.DB) *UsageRecordRepository {
	return &UsageRecordRepository{db: db}
}

// Create creates a new usage record
func (r *UsageRecordRepository) Create(record *entity.UsageRecord) error {
	record.ID = uuid.New().String()

	query := `
		INSERT INTO usage_records (id, user_id, group_id, provider_id, model, prompt_tokens, completion_tokens, cost, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, record.ID, record.UserID, record.GroupID, record.ProviderID, record.Model,
		record.PromptTokens, record.CompletionTokens, record.Cost, record.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	return nil
}

// GetByID retrieves a usage record by ID
func (r *UsageRecordRepository) GetByID(id string) (*entity.UsageRecord, error) {
	query := `
		SELECT id, user_id, group_id, provider_id, model, prompt_tokens, completion_tokens, cost, timestamp
		FROM usage_records
		WHERE id = ?
	`

	var record entity.UsageRecord
	err := r.db.QueryRow(query, id).Scan(
		&record.ID, &record.UserID, &record.GroupID, &record.ProviderID, &record.Model,
		&record.PromptTokens, &record.CompletionTokens, &record.Cost, &record.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usage record not found")
		}
		return nil, fmt.Errorf("failed to get usage record: %w", err)
	}

	return &record, nil
}

// GetByUserID retrieves usage records for a user
func (r *UsageRecordRepository) GetByUserID(userID string, page, pageSize int) ([]*entity.UsageRecord, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM usage_records WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count usage records: %w", err)
	}

	// Get records
	query := `
		SELECT id, user_id, provider_id, model, prompt_tokens, completion_tokens, cost, timestamp
		FROM usage_records
		WHERE user_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query usage records: %w", err)
	}
	defer rows.Close()

	var records []*entity.UsageRecord
	for rows.Next() {
		var record entity.UsageRecord
		err := rows.Scan(
			&record.ID, &record.UserID, &record.ProviderID, &record.Model,
			&record.PromptTokens, &record.CompletionTokens, &record.Cost, &record.Timestamp,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan usage record: %w", err)
		}
		records = append(records, &record)
	}

	return records, total, nil
}

// GetAggregation retrieves aggregated usage statistics
func (r *UsageRecordRepository) GetAggregation(userID, startDate, endDate, groupBy string) ([]*entity.UsageAggregation, error) {
	query := `
		SELECT
			user_id,
			group_id,
			provider_id,
			model,
			COUNT(*) as total_requests,
			SUM(prompt_tokens) as total_prompt_tokens,
			SUM(completion_tokens) as total_completion_tokens,
			SUM(cost) as total_cost
		FROM usage_records
		WHERE user_id = ? AND timestamp BETWEEN ? AND ?
	`

	if groupBy != "" {
		switch groupBy {
		case "user_id":
			query += " GROUP BY user_id"
		case "group_id":
			query += " GROUP BY group_id"
		case "provider_id":
			query += " GROUP BY provider_id"
		case "model":
			query += " GROUP BY model"
		default:
			query += " GROUP BY provider_id, model"
		}
	} else {
		query += " GROUP BY provider_id, model"
	}

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage aggregation: %w", err)
	}
	defer rows.Close()

	var aggregations []*entity.UsageAggregation
	for rows.Next() {
		var agg entity.UsageAggregation
		err := rows.Scan(
			&agg.UserID, &agg.GroupID, &agg.ProviderID, &agg.Model,
			&agg.TotalRequests, &agg.TotalPromptTokens, &agg.TotalCompletionTokens, &agg.TotalCost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aggregation: %w", err)
		}
		agg.StartDate = startDate
		agg.EndDate = endDate
		aggregations = append(aggregations, &agg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if len(aggregations) == 0 {
		return []*entity.UsageAggregation{{
			UserID:    userID,
			StartDate: startDate,
			EndDate:   endDate,
		}}, nil
	}

	return aggregations, nil
}
