package usersrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
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
		eventName = "repo.users.create"
		query     = `
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

func (a usersRepo) FindByUserName(ctx context.Context, username string) (*entity.User, error) {
	var (
		eventName = "repo.users.find_by_user_name"
		query     = `
		SELECT id, user_name, full_name, role
		FROM users
		WHERE user_name = $1
		`
		args = []interface{}{
			username,
		}
		user entity.User
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Fullname, &user.Role)
	if err != nil {
		log.Println(eventName, err)
		return nil, errbank.TranslateDBError(err)
	}

	return &user, nil
}

func (a usersRepo) UpdateRoleByUserName(ctx context.Context, user entity.User) error {
	var (
		eventName = "repo.users.update_role_by_user_name"
		query     = `
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
		log.Println(eventName, err)
		return errbank.TranslateDBError(err)
	}
	return nil
}
