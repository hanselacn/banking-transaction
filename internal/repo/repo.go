package repo

import (
	"database/sql"

	accountrepo "github.com/hanselacn/banking-transaction/internal/repo/account_repo"
	authorizationrepo "github.com/hanselacn/banking-transaction/internal/repo/authorization_repo"
	transactionrepo "github.com/hanselacn/banking-transaction/internal/repo/transaction_repo"
	usersrepo "github.com/hanselacn/banking-transaction/internal/repo/users_repo"
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
