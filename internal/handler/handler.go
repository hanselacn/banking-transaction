package handler

import "database/sql"

type handler struct {
	UsersHandler UsersHandler
}

func NewHandler(db *sql.DB) handler {
	return handler{UsersHandler: NewUsersHandler(db)}
}
