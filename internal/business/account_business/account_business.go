package accountbusiness

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/consts"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
	"github.com/hanselacn/banking-transaction/repo"
)

type AccountBusiness interface {
	Withdrawal(ctx context.Context, input entity.Withdrawal) error
	Deposit(ctx context.Context, input entity.Deposit) error
	GetAccountBalance(ctx context.Context, username string) (*entity.Account, error)
	InterestPayout(ctx context.Context) ([]entity.Account, error)
	InterestPayoutWorker(ctx context.Context) ([]entity.Account, int, int, error)
	UpdateInterestRate(ctx context.Context, input entity.UpdateInterestRate) error
}

type accountBusiness struct {
	repo repo.Repo
	db   *sql.DB
}

func NewAccountBusiness(db *sql.DB) AccountBusiness {
	return &accountBusiness{
		repo: repo.NewRepositories(db),
		db:   db,
	}
}

func (b *accountBusiness) Withdrawal(ctx context.Context, input entity.Withdrawal) error {
	var (
		eventName        = "business.account.withdrawal"
		transactionInput = entity.Transaction{
			ID:     uuid.New(),
			Type:   consts.TxTypeDEBIT,
			Amount: input.Amount,
			Action: consts.TxActionWITHDRAWAL,
			Status: consts.TxStatusINPROGRESS,
		}
	)

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(eventName, "Rollback", rollbackErr)
			}
		}
	}()

	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		log.Println(eventName, err)
		return err
	}
	err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
	if err != nil {
		log.Println(eventName, err)
		return err
	}
	account, err := b.repo.Account.FindByUserID(ctx, user.ID)
	if err != nil {
		log.Println(eventName, err)
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}
	if input.Amount > account.Balance {
		return errbank.NewErrUnprocessableEntity("insufficient balance")
	}

	// Update account balance
	account.Balance -= input.Amount
	err = b.repo.Account.UpdateBalance(ctx, entity.Account{
		UserID:  user.ID,
		Balance: account.Balance,
	}, tx)
	if err != nil {
		log.Println(eventName, "UpdateBalance", err)
		return err
	}
	// Update transaction status to completed
	err = b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusCOMPLETED, tx)
	if err != nil {
		log.Println(eventName, "UpdateTransactionStatus", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Println(eventName, "commiting transaction error", err)
		return err
	}
	return nil
}

func (b *accountBusiness) Deposit(ctx context.Context, input entity.Deposit) error {
	var (
		eventName        = "business.account.deposit"
		transactionInput = entity.Transaction{
			ID:     uuid.New(),
			Type:   consts.TxTypeCREDIT,
			Amount: input.Amount,
			Action: consts.TxActionDEPOSIT,
			Status: consts.TxStatusINPROGRESS,
		}
	)

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(eventName, "Rollback", rollbackErr)
			}
		}
	}()

	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	account, err := b.repo.Account.FindByUserID(ctx, user.ID)
	if err != nil {
		log.Println(eventName, err)
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}

	// Update account balance
	account.Balance += input.Amount
	err = b.repo.Account.UpdateBalance(ctx, entity.Account{
		UserID:  user.ID,
		Balance: account.Balance,
	}, tx)
	if err != nil {
		b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
		log.Println(eventName, "UpdateBalance", err)
		return err
	}
	// Update transaction status to completed
	err = b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusCOMPLETED, tx)
	if err != nil {
		b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
		log.Println(eventName, "UpdateTransactionStatus", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
		log.Println(eventName, "commiting transaction error", err)
		return err
	}
	return nil
}

func (b *accountBusiness) GetAccountBalance(ctx context.Context, username string) (*entity.Account, error) {
	var (
		eventName = "business.account.get_balance"
	)
	user, err := b.repo.Users.FindByUserName(ctx, username)
	if err != nil {
		log.Println(eventName, err)
		return nil, err
	}
	account, err := b.repo.Account.FindByUserID(ctx, user.ID)
	if err != nil {
		log.Println(eventName, err)
		return nil, err
	}
	return account, nil
}

func (b *accountBusiness) UpdateInterestRate(ctx context.Context, input entity.UpdateInterestRate) error {
	var (
		eventName = "business.account.update_interest_rate"
	)

	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	account, err := b.repo.Account.FindByUserID(ctx, user.ID)
	if err != nil {
		log.Println(eventName, err)
		return err
	}

	account.InterestRate = input.InterestRate
	err = b.repo.Account.UpdateInterestRate(ctx, *account, nil)
	if err != nil {
		log.Println(eventName, err)
		return err
	}
	return nil
}

func (b *accountBusiness) InterestPayout(ctx context.Context) ([]entity.Account, error) {
	var (
		eventName = "business.account.get_balance"
	)

	accounts, err := b.repo.Account.GetListAccount(ctx, nil)
	if err != nil {
		log.Println(eventName, err)
		return nil, err
	}

	for i := range accounts {
		tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		})
		if err != nil {
			log.Println(eventName, "BeginTx", err)
			return nil, err
		}

		currentBalance := accounts[i].Balance
		lastPayoutTime := accounts[i].LastInterestPayout.Local()
		currentTime := time.Now()
		durationSinceLastPayout := currentTime.Sub(lastPayoutTime)
		payoutPeriodDays := durationSinceLastPayout.Hours() / 24
		interestAmount := currentBalance * (accounts[i].InterestRate * (payoutPeriodDays / 365))

		transactionInput := entity.Transaction{
			ID:     uuid.New(),
			Type:   consts.TxTypeCREDIT,
			Amount: interestAmount,
			Action: consts.TxActionINTEREST,
			Status: consts.TxStatusINPROGRESS,
		}

		err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
		if err != nil {
			log.Println(eventName, "CreateTransaction", err)
			tx.Rollback()
			return nil, err
		}

		accounts[i].Balance += interestAmount

		err = b.repo.Account.PayoutInterest(ctx, accounts[i], tx)
		if err != nil {
			log.Println(eventName, "PayoutInterest", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println(eventName, "Rollback error", rbErr)
			}
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			continue
		}

		err = b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusCOMPLETED, tx)
		if err != nil {
			log.Println(eventName, "UpdateTransactionStatus", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println(eventName, "Rollback error", rbErr)
			}
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.Println(eventName, "Commit error", err)
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			return nil, err
		}
	}

	return accounts, nil
}

func (b *accountBusiness) InterestPayoutWorker(ctx context.Context) ([]entity.Account, int, int, error) {
	var (
		eventName    = "business.account.get_balance"
		totalSuccess = 0
	)

	accounts, err := b.repo.Account.GetListAccount(ctx, nil)
	if err != nil {
		log.Println(eventName, err)
		return nil, 0, 0, err
	}
	totalData := len(accounts)
	for i := range accounts {
		tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		})
		if err != nil {
			log.Println(eventName, "BeginTx", err)
			return nil, totalData, totalSuccess, err
		}

		currentBalance := accounts[i].Balance
		lastPayoutTime := accounts[i].LastInterestPayout.Local()
		currentTime := time.Now()
		durationSinceLastPayout := currentTime.Sub(lastPayoutTime)
		payoutPeriodDays := durationSinceLastPayout.Hours() / 24
		interestAmount := currentBalance * (accounts[i].InterestRate * (payoutPeriodDays / 365))

		transactionInput := entity.Transaction{
			ID:     uuid.New(),
			Type:   consts.TxTypeCREDIT,
			Amount: interestAmount,
			Action: consts.TxActionINTEREST,
			Status: consts.TxStatusINPROGRESS,
		}

		err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
		if err != nil {
			log.Println(eventName, "CreateTransaction", err)
			tx.Rollback()
			return nil, totalData, totalSuccess, err
		}

		accounts[i].Balance += interestAmount

		err = b.repo.Account.PayoutInterest(ctx, accounts[i], tx)
		if err != nil {
			log.Println(eventName, "PayoutInterest", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println(eventName, "Rollback error", rbErr)
			}
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			continue
		}

		err = b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusCOMPLETED, tx)
		if err != nil {
			log.Println(eventName, "UpdateTransactionStatus", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println(eventName, "Rollback error", rbErr)
			}
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			return nil, totalData, totalSuccess, err
		}

		err = tx.Commit()
		if err != nil {
			log.Println(eventName, "Commit error", err)
			b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID, consts.TxStatusFAILED, nil)
			return nil, totalData, totalSuccess, err
		}
		totalSuccess++
	}
	return accounts, totalData, totalSuccess, nil
}
