package handler

import "database/sql"

type handler struct {
	UsersHandler   UsersHandler
	AccountHandler AccountHandler
}

func NewHandler(db *sql.DB) handler {
	return handler{
		UsersHandler:   NewUsersHandler(db),
		AccountHandler: NewAccountHandler(db),
	}
}
