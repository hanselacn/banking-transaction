package usersbusiness

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"math/big"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hanselacn/banking-transaction/internal/consts"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/repo"
)

type UsersBusiness interface {
	CreateUser(ctx context.Context, input entity.CreateUserInput) (*entity.User, error)
	UpdateRoleByUserName(ctx context.Context, input entity.User) error
	GetUserDetail(ctx context.Context, username string) (*entity.User, error)
}

type usersBusiness struct {
	repo repo.Repo
	db   *sql.DB
}

func NewUsersBusiness(db *sql.DB) UsersBusiness {
	return &usersBusiness{
		repo: repo.NewRepositories(db),
		db:   db,
	}
}

func (b *usersBusiness) CreateUser(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	var (
		eventName = "business.users.create_user"
		user      = entity.User{
			ID:       uuid.New(),
			Username: input.Username,
			Fullname: input.Fullname,
			Role:     consts.RoleCustomer,
		}
	)
	interestEnv := os.Getenv("DEFAULT_INTEREST_RATE")
	defaultInterestRate, err := strconv.ParseFloat(interestEnv, 64)
	if err != nil {
		defaultInterestRate = 0
	}

	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		log.Println(eventName, err)
		return nil, err
	}

	if err := b.repo.Users.Create(ctx, user, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(eventName, err)
			return nil, err
		}
		log.Println(eventName, err)
		return nil, err
	}

	if err := b.repo.Authorization.Create(ctx, entity.Authorization{
		ID:       uuid.New(),
		UserID:   user.ID,
		Password: input.Password,
	}, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(eventName, err)
			return nil, err
		}
		log.Println(eventName, err)
		return nil, err
	}

	num := big.NewInt(0).SetBytes(user.ID[:])
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := big.NewInt(0).Rand(r, num)

	// Format the random number as a string with leading zeros
	randomNumberString := fmt.Sprintf("%036d", randomNumber)
	if err = b.repo.Account.Create(ctx, entity.Account{
		ID:            uuid.New(),
		UserID:        user.ID,
		AccountNumber: randomNumberString,
		Balance:       0,
		InterestRate:  defaultInterestRate,
	}, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(eventName, err)
			return nil, err
		}
		log.Println(eventName, err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(eventName, err)
			return nil, err
		}
		log.Println(eventName, err)
		return nil, err
	}

	return &user, nil
}

func (b *usersBusiness) UpdateRoleByUserName(ctx context.Context, input entity.User) error {
	var (
		eventName = "business.users.update_role"
	)
	if err := b.repo.Users.UpdateRoleByUserName(ctx, input); err != nil {
		log.Println(eventName, err)
		return err
	}
	return nil
}

func (b *usersBusiness) GetUserDetail(ctx context.Context, username string) (*entity.User, error) {
	user, err := b.repo.Users.FindByUserName(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
