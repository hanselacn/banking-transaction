package transactionrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/entity"
)

type TransactionRepo interface {
	CreateTransaction(ctx context.Context, input entity.Transaction) error
	UpdateTransactionStatus(ctx context.Context, trID uuid.UUID, status string, tx *sql.Tx) error
}

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return transactionRepo{db: db}
}

// CreateTransaction implements TransactionRepo.
func (t transactionRepo) CreateTransaction(ctx context.Context, input entity.Transaction) error {
	var (
		query = `
		INSERT INTO transactions (
		id,
		type,
		amount,
		action,
		status,
		updated_at,
		created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
		args = []interface{}{
			input.ID,
			input.Type,
			input.Amount,
			input.Action,
			input.Status,
			time.Now(),
			time.Now(),
		}
	)

	_, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTransaction implements TransactionRepo.
func (t transactionRepo) UpdateTransactionStatus(ctx context.Context, trID uuid.UUID, status string, tx *sql.Tx) error {
	var (
		query = `
		UPDATE transactions
		SET status = $1
		WHERE id = $2
		`
		args = []interface{}{
			status,
			trID,
		}
	)
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	_, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
