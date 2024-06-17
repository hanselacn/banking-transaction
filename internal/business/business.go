package business

import (
	"database/sql"

	authorizationbusiness "github.com/hanselacn/banking-transaction/internal/business/authorization_business"
	usersbusiness "github.com/hanselacn/banking-transaction/internal/business/users_business"
)

type Business struct {
	UserBusiness          usersbusiness.UsersBusiness
	AuthorizationBusiness authorizationbusiness.AuthorizationBusiness
}

func NewBusiness(db *sql.DB) Business {
	return Business{
		UserBusiness:          usersbusiness.NewUsersBusiness(db),
		AuthorizationBusiness: authorizationbusiness.NewAuthorizationBusiness(db),
	}
}
