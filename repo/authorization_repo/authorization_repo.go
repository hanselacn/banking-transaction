package authorizationrepo

import (
	"context"

	"github.com/hanselacn/banking-transaction/internal/entity"
)

type AuthorizationRepositories interface {
	FindByAccountNumber(ctx context.Context, accountNumber string)
	Udpate(ctx context.Context, authorization entity.Authorization)
	Create(ctx context.Context, authorization entity.Authorization)
	DeleteByID(ctx context.Context, accountNumber string)
}

type authorizationrepo struct {
}

// Create implements AuthorizationRepositories.
func (a authorizationrepo) Create(ctx context.Context, authorization entity.Authorization) {
	panic("unimplemented")
}

// DeleteByID implements AuthorizationRepositories.
func (a authorizationrepo) DeleteByID(ctx context.Context, accountNumber string) {
	panic("unimplemented")
}

// FindByAccountNumber implements AuthorizationRepositories.
func (a authorizationrepo) FindByAccountNumber(ctx context.Context, accountNumber string) {
	panic("unimplemented")
}

// Udpate implements AuthorizationRepositories.
func (a authorizationrepo) Udpate(ctx context.Context, authorization entity.Authorization) {
	panic("unimplemented")
}

func NewAccountRepositories() AuthorizationRepositories {
	return authorizationrepo{}
}
