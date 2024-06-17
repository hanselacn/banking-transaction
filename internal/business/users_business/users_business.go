package usersbusiness

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/consts"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/repo"
)

type UsersBusiness interface {
}

type usersBusiness struct {
	repo repo.Repo
	db   *sql.DB
}

func NewUsersBusiness(db *sql.DB) UsersBusiness {
	return usersBusiness{
		repo: repo.NewRepositories(db),
		db:   db,
	}
}

func (b *usersBusiness) CreateUser(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	var (
		user = entity.User{
			ID:       uuid.New(),
			Username: input.Username,
			Fullname: input.Fullname,
			Role:     consts.RoleCustomer,
		}
	)

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	if err := b.repo.Users.Create(ctx, user, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := b.repo.Authorization.Create(ctx, entity.Authorization{
		ID:       uuid.New(),
		UserID:   user.ID,
		Password: input.Password,
	}, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (b *usersBusiness) UpdateRoleByUserName(ctx context.Context, input entity.User) error {
	if err := b.repo.Users.UpdateRoleByUserName(ctx, input); err != nil {
		return err
	}
	return nil
}
