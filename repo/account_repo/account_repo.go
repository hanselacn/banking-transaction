package accountrepo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/entity"
)

type AccountRepositories interface {
	FindByUserID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	Create(ctx context.Context, account entity.Account) error
	UpdateBalance(ctx context.Context, account entity.Account, tx *sql.Tx) error
	UpdateInterestRate(ctx context.Context, account entity.Account) error
}

type accountRepo struct {
	db *sql.DB
}

// Create implements AccountRepositories.
func (a accountRepo) Create(ctx context.Context, account entity.Account) error {
	var (
		query = `
		INSERT INTO account (
		id,
		user_id,
		account_number,
		balance,
		interest_rate
		)
		VALUES ($1,$2,$3,$4,$5)
	`
		args = []interface{}{
			account.ID,
			account.UserID,
			account.AccountNumber,
			account.Balance,
			account.InterestRate,
		}
	)

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (a accountRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var (
		query = `
		SELECT id, user_id, account_number, balance, interest_rate
		FROM account
		WHERE user_id = $1
		`
		args = []interface{}{
			id,
		}
		account entity.Account
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a accountRepo) UpdateBalance(ctx context.Context, account entity.Account, tx *sql.Tx) error {
	var (
		query = `
		UPDATE account
		SET balance = $1
		WHERE user_id = $2
		`
		args = []interface{}{
			account.Balance,
			account.UserID,
		}
	)

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	} else {
		_, err := a.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a accountRepo) UpdateInterestRate(ctx context.Context, account entity.Account) error {
	var (
		query = `
		UPDATE account
		SET interest_rate = $1
		WHERE user_id = $1
		`
		args = []interface{}{
			account.InterestRate,
			account.UserID,
		}
	)

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewAccountRepositories(db *sql.DB) AccountRepositories {
	return accountRepo{db: db}
}
