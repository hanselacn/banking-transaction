package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"github.com/hanselacn/banking-transaction/internal/business"
	"github.com/hanselacn/banking-transaction/internal/consts"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/middleware"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
	"github.com/hanselacn/banking-transaction/internal/pkg/response"
	"github.com/hanselacn/banking-transaction/internal/pkg/rule"
	"github.com/pkg/errors"
)

func MountAccountHandler(r *mux.Router, h handler, m middleware.Middleware) {
	r.Handle("/banking-transaction/account/withdrawal", m.AuthenticationMiddleware((http.HandlerFunc(h.AccountHandler.Withdrawal)), []string{consts.RoleSuperAdmin, consts.RoleAdmin, consts.RoleCustomer}...)).Methods("POST")
	r.Handle("/banking-transaction/account/deposit", m.AuthenticationMiddleware((http.HandlerFunc(h.AccountHandler.Deposit)), []string{consts.RoleSuperAdmin, consts.RoleAdmin, consts.RoleCustomer}...)).Methods("POST")
	r.Handle("/banking-transaction/account/balance/{user_name}", m.AuthenticationMiddleware((http.HandlerFunc(h.AccountHandler.GetAccountBalance)), []string{consts.RoleSuperAdmin, consts.RoleAdmin, consts.RoleCustomer}...)).Methods("GET")
	r.Handle("/banking-transaction/account/interest/payout", m.AuthenticationMiddleware((http.HandlerFunc(h.AccountHandler.InterestPayout)), []string{consts.RoleSuperAdmin, consts.RoleAdmin}...)).Methods("POST")
	r.Handle("/banking-transaction/account/interest/update", m.AuthenticationMiddleware((http.HandlerFunc(h.AccountHandler.UpdateInterestRate)), []string{consts.RoleSuperAdmin, consts.RoleAdmin}...)).Methods("PUT")
}

type AccountHandler struct {
	business business.Business
}

func NewAccountHandler(db *sql.DB) AccountHandler {
	return AccountHandler{business: business.NewBusiness(db)}
}

func (h *AccountHandler) Withdrawal(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.account.withdrawal"
		payload   entity.Withdrawal
	)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "request body malformed", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := payload.Validate(); err != nil {
		response.JsonResponse(w, "withdrawal error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if payload.Amount > 1000000000000 || payload.Amount < 1 {
		response.JsonResponse(w, "withdrawal error", nil, "ammount must between 1-1000000000000", http.StatusUnprocessableEntity)
		return
	}

	ctxUserName := ctx.Value(middleware.CtxValueUserName)
	if ctxUserName != payload.Username {
		response.JsonResponse(w, "Forbidden", nil, "You Have to Access your own Account", http.StatusForbidden)
		return
	}

	err = h.business.AccountBusiness.Withdrawal(ctx, payload)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		log.Println(eventName, err)
		causer := errors.Cause(err)
		switch causer.(type) {
		case errbank.ErrConflict:
			statusCode = http.StatusConflict
		case errbank.ErrNotFound:
			statusCode = http.StatusNotFound
		case errbank.ErrUnprocessableEntity:
			statusCode = http.StatusUnprocessableEntity
		case errbank.ErrForbidden:
			statusCode = http.StatusForbidden
		case errbank.ErrTooManyRequest:
			statusCode = http.StatusTooManyRequests
		}
		response.JsonResponse(w, "withdrawal error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success withdrawal", nil, nil, http.StatusOK)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.account.deposit"
		payload   entity.Deposit
	)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "request body malformed", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := payload.Validate(); err != nil {
		response.JsonResponse(w, "deposit error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if payload.Amount > 1000000000000 || payload.Amount < 1 {
		response.JsonResponse(w, "deposit error", nil, "ammount must between 1-1000000000000", http.StatusUnprocessableEntity)
		return
	}

	ctxUserName := ctx.Value(middleware.CtxValueUserName)
	if ctxUserName != payload.Username {
		response.JsonResponse(w, "Forbidden", nil, "You Have to Access your own Account", http.StatusForbidden)
		return
	}
	err = h.business.AccountBusiness.Deposit(ctx, payload)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		log.Println(eventName, err)
		causer := errors.Cause(err)
		switch causer.(type) {
		case errbank.ErrConflict:
			statusCode = http.StatusConflict
		case errbank.ErrNotFound:
			statusCode = http.StatusNotFound
		case errbank.ErrUnprocessableEntity:
			statusCode = http.StatusUnprocessableEntity
		case errbank.ErrForbidden:
			statusCode = http.StatusForbidden
		case errbank.ErrTooManyRequest:
			statusCode = http.StatusTooManyRequests
		}
		response.JsonResponse(w, "deposit error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success deposit", nil, nil, http.StatusOK)
}

func (h *AccountHandler) UpdateInterestRate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.account.update_interest_rate"
		payload   entity.UpdateInterestRate
	)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "request body malformed", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := payload.Validate(); err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "update interest rate error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = h.business.AccountBusiness.UpdateInterestRate(ctx, payload)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		log.Println(eventName, err)
		causer := errors.Cause(err)
		switch causer.(type) {
		case errbank.ErrConflict:
			statusCode = http.StatusConflict
		case errbank.ErrNotFound:
			statusCode = http.StatusNotFound
		case errbank.ErrUnprocessableEntity:
			statusCode = http.StatusUnprocessableEntity
		case errbank.ErrForbidden:
			statusCode = http.StatusForbidden
		case errbank.ErrTooManyRequest:
			statusCode = http.StatusTooManyRequests
		}
		response.JsonResponse(w, "deposit error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success deposit", nil, nil, http.StatusOK)
}

func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.account.get_balance"
		pathVar   = mux.Vars(r)
		username  = pathVar["user_name"]
	)

	ctxUserName := ctx.Value(middleware.CtxValueUserName)
	if ctxUserName != username {
		response.JsonResponse(w, "Forbidden", nil, "You Have to Access your own Account", http.StatusForbidden)
		return
	}

	if err := validation.Validate(username, rule.UserNameRule); err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "get account balance error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	account, err := h.business.AccountBusiness.GetAccountBalance(ctx, username)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		log.Println(eventName, err)
		causer := errors.Cause(err)
		switch causer.(type) {
		case errbank.ErrConflict:
			statusCode = http.StatusConflict
		case errbank.ErrNotFound:
			statusCode = http.StatusNotFound
		case errbank.ErrUnprocessableEntity:
			statusCode = http.StatusUnprocessableEntity
		case errbank.ErrForbidden:
			statusCode = http.StatusForbidden
		case errbank.ErrTooManyRequest:
			statusCode = http.StatusTooManyRequests
		}
		response.JsonResponse(w, "get account balance error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success get account balance", account, nil, http.StatusOK)
}

func (h *AccountHandler) InterestPayout(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.account.interest_payout"
	)

	role := ctx.Value(middleware.CtxValueRole)
	switch role {
	case "customer":
		response.JsonResponse(w, "Forbidden", nil, "You Cannot Access This Feature", http.StatusForbidden)
		return
	}

	account, err := h.business.AccountBusiness.InterestPayout(ctx)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		log.Println(eventName, err)
		causer := errors.Cause(err)
		switch causer.(type) {
		case errbank.ErrConflict:
			statusCode = http.StatusConflict
		case errbank.ErrNotFound:
			statusCode = http.StatusNotFound
		case errbank.ErrUnprocessableEntity:
			statusCode = http.StatusUnprocessableEntity
		case errbank.ErrForbidden:
			statusCode = http.StatusForbidden
		case errbank.ErrTooManyRequest:
			statusCode = http.StatusTooManyRequests
		}
		response.JsonResponse(w, "interest compounding error error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success interest compounding", account, nil, http.StatusOK)
}
