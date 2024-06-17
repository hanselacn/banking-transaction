package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/pkg/rule"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"user_name"`
	Fullname string    `json:"full_name"`
	Role     string    `json:"role"`
}

func (s *User) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Username, validation.Required, rule.UserNameRule),
		validation.Field(&s.Role, validation.In("admin", "super_admin", "customer").Error("invalid role input")),
	)
}

type CreateUserInput struct {
	Username string `json:"user_name"`
	Fullname string `json:"full_name"`
	Password string `json:"password"`
}

func (s *CreateUserInput) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Username, validation.Required, rule.UserNameRule),
		validation.Field(&s.Fullname, rule.AlphabetNumericSpaceCharRule),
		validation.Field(&s.Password, rule.SpecialCharRegexRule),
		validation.Field(&s.Password, rule.LengthRegexRule),
		validation.Field(&s.Password, rule.LowercaseRegexRule),
		validation.Field(&s.Password, rule.UppercaseRegexRule),
		validation.Field(&s.Password, rule.DigitRegexRule),
	)
}

type Withdrawal struct {
	Username string  `json:"user_name"`
	Amount   float64 `json:"amount"`
}

func (s *Withdrawal) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Username, validation.Required, rule.UserNameRule),
	)
}

type Deposit struct {
	Username string  `json:"user_name"`
	Amount   float64 `json:"amount"`
}

func (s *Deposit) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Username, validation.Required, rule.UserNameRule),
	)
}

type UpdateInterestRate struct {
	Username     string  `json:"user_name"`
	InterestRate float64 `json:"interest_rate"`
}

func (s *UpdateInterestRate) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Username, validation.Required, rule.UserNameRule),
		validation.Field(&s.InterestRate, rule.InterestRateRule),
	)
}

type Account struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	AccountNumber      string    `json:"account_number"`
	Balance            float64   `json:"balance"`
	InterestRate       float64   `json:"interest_rate"`
	CreatedAt          time.Time `json:"created_at"`
	LastInterestPayout time.Time `json:"last_interest_payout,omitempty"`
}

type AccountPresentation struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	AccountNumber      string    `json:"account_number"`
	Balance            string    `json:"balance"`
	InterestRate       string    `json:"interest_rate"`
	CreatedAt          time.Time `json:"created_at"`
	LastInterestPayout time.Time `json:"last_interest_payout,omitempty"`
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
