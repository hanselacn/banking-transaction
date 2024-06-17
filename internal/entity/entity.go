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

type CreateUserInput struct {
	Username string
	Fullname string
	Password string
}

type Withdrawal struct {
	Username string
	Amount   float64
}

type Deposit struct {
	Username string
	Amount   float64
}

type Account struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	AccountNumber string
	Balance       float64
	InterestRate  float64
}

type ChangePasswordInput struct {
	Username string
	Password string
}

type Authorization struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Password string
}

type Transaction struct {
	ID     uuid.UUID
	Type   string
	Amount float64
	Action string
	Status string
}
