package repo

import (
	"database/sql"

	accountrepo "github.com/hanselacn/banking-transaction/repo/account_repo"
	authorizationrepo "github.com/hanselacn/banking-transaction/repo/authorization_repo"
)

type Repo struct {
	Account       accountrepo.AccountRepositories
	Authorization authorizationrepo.AuthorizationRepositories
}

func NewRepositories(db *sql.DB) Repo {
	return Repo{Account: accountrepo.NewAccountRepositories(db)}
}
