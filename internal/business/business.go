package business

import (
	"database/sql"

	accountbusiness "github.com/hanselacn/banking-transaction/internal/business/account_business"
	authorizationbusiness "github.com/hanselacn/banking-transaction/internal/business/authorization_business"
	usersbusiness "github.com/hanselacn/banking-transaction/internal/business/users_business"
)

type Business struct {
	AccountBusiness       accountbusiness.AccountBusiness
	UserBusiness          usersbusiness.UsersBusiness
	AuthorizationBusiness authorizationbusiness.AuthorizationBusiness
}

func NewBusiness(db *sql.DB) Business {
	return Business{
		AccountBusiness:       accountbusiness.NewAccountBusiness(db),
		UserBusiness:          usersbusiness.NewUsersBusiness(db),
		AuthorizationBusiness: authorizationbusiness.NewAuthorizationBusiness(db),
	}
}
