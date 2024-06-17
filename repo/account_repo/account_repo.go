package accountrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/cryptox"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
)

type AccountRepositories interface {
	FindByUserID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	Create(ctx context.Context, account entity.Account, tx *sql.Tx) error
	UpdateBalance(ctx context.Context, account entity.Account, tx *sql.Tx) error
	UpdateInterestRate(ctx context.Context, account entity.Account, tx *sql.Tx) error
}

type accountRepo struct {
	db *sql.DB
}

// Create implements AccountRepositories.
func (a accountRepo) Create(ctx context.Context, account entity.Account, tx *sql.Tx) error {
	var (
		eventName = "repo.account.create_account"
		query     = `
		INSERT INTO accounts (
		id,
		user_id,
		account_number,
		balance,
		interest_rate
		)
		VALUES ($1,$2,$3,$4,$5)
	`
	)

	balance := fmt.Sprintf("%.2f", account.Balance)
	interestRate := fmt.Sprintf("%.2f", account.InterestRate)

	encryptedBalance, err := cryptox.EncryptAES(balance, os.Getenv("AES_KEY"))
	if err != nil {
		return err
	}

	encryptedInterestRate, err := cryptox.EncryptAES(interestRate, os.Getenv("AES_KEY"))
	if err != nil {
		return err
	}

	args := []interface{}{
		account.ID,
		account.UserID,
		account.AccountNumber,
		encryptedBalance,
		encryptedInterestRate,
	}

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}
	} else {
		_, err := a.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}
	}

	return nil
}

func (a accountRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var (
		eventName = "repo.account.find_user_by_id"
		query     = `
		SELECT id, user_id, account_number, balance, interest_rate
		FROM accounts
		WHERE user_id = $1
		`
		args = []interface{}{
			id,
		}
		accountPrs entity.AccountPresentation
	)

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&accountPrs.ID, &accountPrs.UserID, &accountPrs.AccountNumber, &accountPrs.Balance, &accountPrs.InterestRate)
	if err != nil {
		log.Println(eventName, err)
		return nil, errbank.TranslateDBError(err)
	}

	accountPrs.Balance, err = cryptox.DecryptAES(accountPrs.Balance, os.Getenv("AES_KEY"))
	if err != nil {
		return nil, err
	}
	accountPrs.InterestRate, err = cryptox.DecryptAES(accountPrs.InterestRate, os.Getenv("AES_KEY"))
	if err != nil {
		return nil, err
	}

	balance, err := strconv.ParseFloat(accountPrs.Balance, 64)
	if err != nil {
		return nil, err
	}

	interest, err := strconv.ParseFloat(accountPrs.InterestRate, 64)
	if err != nil {
		return nil, err
	}

	return &entity.Account{
		ID:            accountPrs.ID,
		UserID:        accountPrs.UserID,
		AccountNumber: accountPrs.AccountNumber,
		Balance:       balance,
		InterestRate:  interest,
	}, nil
}
func (a accountRepo) UpdateBalance(ctx context.Context, account entity.Account, tx *sql.Tx) error {
	var (
		eventName = "repo.account.update_balance"
		query     = `
		UPDATE accounts
		SET balance = $1
		WHERE user_id = $2
		`
	)

	// Convert balance to string and encrypt
	balance := fmt.Sprintf("%.2f", account.Balance)
	encryptedBalance, err := cryptox.EncryptAES(balance, os.Getenv("AES_KEY"))
	if err != nil {
		log.Println(eventName, "EncryptAES", err)
		return err
	}

	// Prepare arguments
	args := []interface{}{
		encryptedBalance,
		account.UserID,
	}

	// Execute update query within the transaction if provided
	if tx != nil {
		log.Println(eventName, "Executing within transaction")
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, "ExecContext (tx)", err)
			return errbank.TranslateDBError(err)
		}
	} else {
		log.Println(eventName, "Executing without transaction")
		_, err := a.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, "ExecContext (no tx)", err)
			return errbank.TranslateDBError(err)
		}
	}

	log.Println(eventName, "Update successful")
	return nil
}

func (a accountRepo) UpdateInterestRate(ctx context.Context, account entity.Account, tx *sql.Tx) error {
	var (
		eventName = "repo.account.update_interest_rate"
		query     = `
		UPDATE accounts
		SET interest_rate = $1
		WHERE user_id = $1
		`
	)
	interest := fmt.Sprintf("%.2f", account.InterestRate)
	encryptedInterestRate, err := cryptox.EncryptAES(interest, os.Getenv("AES_KEY"))
	if err != nil {
		return err
	}

	args := []interface{}{
		encryptedInterestRate,
		account.UserID,
	}

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}
	} else {
		_, err = a.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(eventName, err)
			return errbank.TranslateDBError(err)
		}
	}
	return nil
}

func NewAccountRepositories(db *sql.DB) AccountRepositories {
	return accountRepo{db: db}
}
