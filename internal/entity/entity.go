package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"user_name"`
	Fullname string    `json:"full_name"`
	Role     string    `json:"role"`
}

type CreateUserInput struct {
	Username string `json:"user_name"`
	Fullname string `json:"full_name"`
	Password string `json:"password"`
}

type Withdrawal struct {
	Username string  `json:"user_name"`
	Amount   float64 `json:"amount"`
}

type Deposit struct {
	Username string  `json:"user_name"`
	Amount   float64 `json:"amount"`
}

type Account struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	InterestRate  float64   `json:"interest_rate"`
}

type AccountPresentation struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	Balance       string    `json:"balance"`
	InterestRate  string    `json:"interest_rate"`
}

type ChangePasswordInput struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type Authorization struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Type   string    `json:"type"`
	Amount float64   `json:"amount"`
	Action string    `json:"action"`
	Status string    `json:"status"`
}
