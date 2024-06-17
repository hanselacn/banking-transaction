package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hanselacn/banking-transaction/internal/business"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
	"github.com/hanselacn/banking-transaction/internal/pkg/response"
	"github.com/pkg/errors"
)

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
