package authorizationrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
)

type AuthorizationRepositories interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Authorization, error)
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
		eventName = "repo.authorization.create"
		query     = `
		INSERT INTO authorizations (
		id,
		user_id,
		password
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

			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}

	} else {
		_, err := a.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}
	}

	return nil
}

func (a authorizationrepo) UpdatePassword(ctx context.Context, authorization entity.Authorization) error {
	var (
		eventName = "repo.authorization.update_password"
		query     = `
		UPDATE authorizations
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
		log.Println(eventName, err)
		return errbank.TranslateDBError(err)
	}
	return nil
}

func (a authorizationrepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Authorization, error) {
	var (
		eventName = "repo.authorization.find_user_by_id"
		query     = `
		SELECT id, user_id, password
		FROM authorizations
		WHERE user_id = $1
		`
		args = []interface{}{
			userID,
		}
		auth entity.Authorization
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&auth.ID, &auth.UserID, &auth.Password)
	if err != nil {
		log.Println(eventName, err)
		return nil, errbank.TranslateDBError(err)
	}
	return &auth, nil
}
