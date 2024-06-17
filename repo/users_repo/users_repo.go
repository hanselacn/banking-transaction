package usersrepo

import (
	"context"
	"database/sql"

	"github.com/hanselacn/banking-transaction/internal/entity"
)

type UsersRepositories interface {
	FindByUserName(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, user entity.User, tx *sql.Tx) error
	UpdateRoleByUserName(ctx context.Context, user entity.User) error
}

type usersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) UsersRepositories {
	return usersRepo{db: db}
}

func (a usersRepo) Create(ctx context.Context, user entity.User, tx *sql.Tx) error {
	var (
		query = `
		INSERT INTO users (
		id,
		user_name,
		full_name,
		role
		)
		VALUES ($1,$2,$3,$4)
	`
		args = []interface{}{
			user.ID,
			user.Username,
			user.Fullname,
			user.Role,
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

func (a usersRepo) FindByUserName(ctx context.Context, username string) (*entity.User, error) {
	var (
		query = `
		SELECT id, user_name, full_name, role
		FROM users
		WHERE user_name = $1
		`
		args = []interface{}{
			username,
		}
		user entity.User
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a usersRepo) UpdateRoleByUserName(ctx context.Context, user entity.User) error {
	var (
		query = `
		UPDATE users
		SET role = $1
		WHERE user_name = $2
		`
		args = []interface{}{
			user.Role,
			user.Username,
		}
	)

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
