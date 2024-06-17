package authorizationbusiness

import (
	"context"
	"database/sql"

	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/repo"
)

type AuthorizationBusiness interface {
}

type authorizationBusiness struct {
	repo repo.Repo
}

func NewAuthorizationBusiness(db *sql.DB) AuthorizationBusiness {
	return authorizationBusiness{repo: repo.NewRepositories(db)}
}

func (b *authorizationBusiness) ChangePassword(ctx context.Context, input entity.ChangePasswordInput) error {
	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		return err
	}

	if err := b.repo.Authorization.UpdatePassword(ctx, entity.Authorization{
		UserID:   user.ID,
		Password: input.Password,
	}); err != nil {
		return err
	}
	return nil
}
