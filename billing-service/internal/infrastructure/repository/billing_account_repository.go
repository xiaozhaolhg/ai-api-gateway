package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// BillingAccountRepository implements the billing account repository interface
type BillingAccountRepository struct {
	db *sql.DB
}

// NewBillingAccountRepository creates a new billing account repository
func NewBillingAccountRepository(db *sql.DB) *BillingAccountRepository {
	return &BillingAccountRepository{db: db}
}

// Create creates a new billing account
func (r *BillingAccountRepository) Create(account *entity.BillingAccount) error {
	account.ID = uuid.New().String()

	query := `
		INSERT INTO billing_accounts (id, user_id, balance, currency, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, account.ID, account.UserID, account.Balance,
		account.Currency, account.Status, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create billing account: %w", err)
	}

	return nil
}

// GetByID retrieves a billing account by ID
func (r *BillingAccountRepository) GetByID(id string) (*entity.BillingAccount, error) {
	query := `
		SELECT id, user_id, balance, currency, status, created_at, updated_at
		FROM billing_accounts
		WHERE id = ?
	`

	var account entity.BillingAccount
	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.UserID, &account.Balance,
		&account.Currency, &account.Status, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("billing account not found")
		}
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	return &account, nil
}

// GetByUserID retrieves a billing account by user ID
func (r *BillingAccountRepository) GetByUserID(userID string) (*entity.BillingAccount, error) {
	query := `
		SELECT id, user_id, balance, currency, status, created_at, updated_at
		FROM billing_accounts
		WHERE user_id = ?
	`

	var account entity.BillingAccount
	err := r.db.QueryRow(query, userID).Scan(
		&account.ID, &account.UserID, &account.Balance,
		&account.Currency, &account.Status, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("billing account not found")
		}
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	return &account, nil
}

// Update updates an existing billing account
func (r *BillingAccountRepository) Update(account *entity.BillingAccount) error {
	query := `
		UPDATE billing_accounts
		SET balance = ?, currency = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, account.Balance, account.Currency,
		account.Status, account.UpdatedAt, account.ID)
	if err != nil {
		return fmt.Errorf("failed to update billing account: %w", err)
	}

	return nil
}

// Delete deletes a billing account
func (r *BillingAccountRepository) Delete(id string) error {
	query := `DELETE FROM billing_accounts WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete billing account: %w", err)
	}

	return nil
}
