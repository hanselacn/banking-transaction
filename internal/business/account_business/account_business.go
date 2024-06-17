package accountbusiness

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/consts"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
	"github.com/hanselacn/banking-transaction/repo"
)

type AccountBusiness interface {
}

type accountBusiness struct {
	repo repo.Repo
	db   *sql.DB
}

func NewAccountBusiness(db *sql.DB) AccountBusiness {
	return accountBusiness{
		repo: repo.NewRepositories(db),
		db:   db,
	}
}

func (b *accountBusiness) Withdrawal(ctx context.Context, input entity.Withdrawal) error {
	transactionInput := entity.Transaction{
		ID:     uuid.New(),
		Type:   consts.TxTypeDEBIT,
		Amount: input.Amount,
		Action: consts.TxActionWITHDRAWAL,
		Status: consts.TxStatusINPROGRESS,
	}
	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		return err
	}
	err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
	if err != nil {
		return err
	}

	account, err := b.repo.Account.FindByUserID(ctx, user.ID.String())
	if err != nil {
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}

	if input.Amount > account.Balance {
		return errbank.NewErrUnprocessableEntity("insufficient balance")
	}

	account.Balance -= input.Amount

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	err = b.repo.Account.UpdateBalance(ctx, entity.Account{
		UserID:  user.ID,
		Balance: account.Balance,
	}, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}

	if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusCOMPLETED, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *accountBusiness) Deposit(ctx context.Context, input entity.Deposit) error {
	transactionInput := entity.Transaction{
		ID:     uuid.New(),
		Type:   consts.TxTypeCREDIT,
		Amount: input.Amount,
		Action: consts.TxActionDEPOSIT,
		Status: consts.TxStatusINPROGRESS,
	}
	user, err := b.repo.Users.FindByUserName(ctx, input.Username)
	if err != nil {
		return err
	}

	err = b.repo.Transaction.CreateTransaction(ctx, transactionInput)
	if err != nil {
		return err
	}

	account, err := b.repo.Account.FindByUserID(ctx, user.ID.String())
	if err != nil {
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}

	account.Balance += input.Amount

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	err = b.repo.Account.UpdateBalance(ctx, entity.Account{
		UserID:  user.ID,
		Balance: account.Balance,
	}, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
		return err
	}

	if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusCOMPLETED, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		if err := b.repo.Transaction.UpdateTransactionStatus(ctx, transactionInput.ID.String(), consts.TxStatusFAILED, nil); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
