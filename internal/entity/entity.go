package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Fullname string
	Role     string
}

type Account struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	AccountNumber string
	Balance       float64
	InterestRate  float64
}

type Authorization struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Password string
}
