package authorizationrepo

import (
	"context"
	"database/sql"

	"github.com/hanselacn/banking-transaction/internal/entity"
)

type AuthorizationRepositories interface {
	FindByUserID(ctx context.Context, userID string) (*entity.Authorization, error)
	UpdatePassword(ctx context.Context, authorization entity.Authorization) error
	Create(ctx context.Context, authorization entity.Authorization, tx *sql.Tx) error
}

type authorizationrepo struct {
	db *sql.DB
}

func NewAccountRepositories(db *sql.DB) AuthorizationRepositories {
	return authorizationrepo{db: db}
}

// Create implements AuthorizationRepositories.
func (a authorizationrepo) Create(ctx context.Context, authorization entity.Authorization, tx *sql.Tx) error {
	var (
		query = `
		INSERT INTO authorization (
		id,
		user_id,
		pasword
		)
		VALUES ($1,$2,$3)
	`
		args = []interface{}{
			authorization.ID,
			authorization.UserID,
			authorization.Password,
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

func (a authorizationrepo) UpdatePassword(ctx context.Context, authorization entity.Authorization) error {
	var (
		query = `
		UPDATE authorization
		SET password = $1
		WHERE user_id = $2
		`
		args = []interface{}{
			authorization.Password,
			authorization.UserID,
		}
	)

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a authorizationrepo) FindByUserID(ctx context.Context, userID string) (*entity.Authorization, error) {
	var (
		query = `
		SELECT id, user_id, password
		FROM authorization
		WHERE user_id = $1
		`
		args = []interface{}{
			userID,
		}
		auth entity.Authorization
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
