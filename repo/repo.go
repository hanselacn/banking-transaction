package repo

import (
	"database/sql"

	accountrepo "github.com/hanselacn/banking-transaction/repo/account_repo"
	authorizationrepo "github.com/hanselacn/banking-transaction/repo/authorization_repo"
	transactionrepo "github.com/hanselacn/banking-transaction/repo/transaction_repo"
	usersrepo "github.com/hanselacn/banking-transaction/repo/users_repo"
)

type Repo struct {
	Account       accountrepo.AccountRepositories
	Authorization authorizationrepo.AuthorizationRepositories
	Users         usersrepo.UsersRepositories
	Transaction   transactionrepo.TransactionRepo
}

func NewRepositories(db *sql.DB) Repo {
	return Repo{
		Account:       accountrepo.NewAccountRepositories(db),
		Authorization: authorizationrepo.NewAccountRepositories(db),
		Users:         usersrepo.NewUsersRepo(db),
		Transaction:   transactionrepo.NewTransactionRepo(db)}
}
