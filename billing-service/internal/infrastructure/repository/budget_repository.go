package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// BudgetRepository implements the budget repository interface
type BudgetRepository struct {
	db *sql.DB
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create creates a new budget
func (r *BudgetRepository) Create(budget *entity.Budget) error {
	budget.ID = uuid.New().String()

	query := `
		INSERT INTO budgets (id, user_id, "limit", period, soft_cap, hard_cap, start_date, end_date, alert_threshold, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, budget.ID, budget.UserID, budget.Limit,
		budget.Period, budget.SoftCap, budget.HardCap,
		budget.StartDate, budget.EndDate, budget.AlertThreshold,
		budget.CreatedAt, budget.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create budget: %w", err)
	}

	return nil
}

// GetByID retrieves a budget by ID
func (r *BudgetRepository) GetByID(id string) (*entity.Budget, error) {
	query := `
		SELECT id, user_id, "limit", period, soft_cap, hard_cap, start_date, end_date, alert_threshold, created_at, updated_at
		FROM budgets
		WHERE id = ?
	`

	var budget entity.Budget
	err := r.db.QueryRow(query, id).Scan(
		&budget.ID, &budget.UserID, &budget.Limit,
		&budget.Period, &budget.SoftCap, &budget.HardCap,
		&budget.StartDate, &budget.EndDate, &budget.AlertThreshold,
		&budget.CreatedAt, &budget.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("budget not found")
		}
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	return &budget, nil
}

// GetByUserID retrieves a budget by user ID
func (r *BudgetRepository) GetByUserID(userID string) (*entity.Budget, error) {
	query := `
		SELECT id, user_id, "limit", period, soft_cap, hard_cap, start_date, end_date, alert_threshold, created_at, updated_at
		FROM budgets
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var budget entity.Budget
	err := r.db.QueryRow(query, userID).Scan(
		&budget.ID, &budget.UserID, &budget.Limit,
		&budget.Period, &budget.SoftCap, &budget.HardCap,
		&budget.StartDate, &budget.EndDate, &budget.AlertThreshold,
		&budget.CreatedAt, &budget.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("budget not found")
		}
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	return &budget, nil
}

// Update updates an existing budget
func (r *BudgetRepository) Update(budget *entity.Budget) error {
	query := `
		UPDATE budgets
		SET limit = ?, period = ?, soft_cap = ?, hard_cap = ?, start_date = ?, end_date = ?, alert_threshold = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, budget.Limit, budget.Period,
		budget.SoftCap, budget.HardCap,
		budget.StartDate, budget.EndDate, budget.AlertThreshold,
		budget.UpdatedAt, budget.ID)
	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	return nil
}

// Delete deletes a budget
func (r *BudgetRepository) Delete(id string) error {
	query := `DELETE FROM budgets WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	return nil
}

// List lists all budgets
func (r *BudgetRepository) List(page, pageSize int) ([]*entity.Budget, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM budgets").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budgets: %w", err)
	}

	// Get budgets
	query := `
		SELECT id, user_id, "limit", period, soft_cap, hard_cap, start_date, end_date, alert_threshold, created_at, updated_at
		FROM budgets
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query budgets: %w", err)
	}
	defer rows.Close()

	var budgets []*entity.Budget
	for rows.Next() {
		var budget entity.Budget
		err := rows.Scan(
			&budget.ID, &budget.UserID, &budget.Limit,
			&budget.Period, &budget.SoftCap, &budget.HardCap,
			&budget.StartDate, &budget.EndDate, &budget.AlertThreshold,
			&budget.CreatedAt, &budget.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, &budget)
	}

	return budgets, total, nil
}
