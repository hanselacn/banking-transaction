package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hanselacn/banking-transaction/internal/business"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/pkg/errbank"
	"github.com/hanselacn/banking-transaction/internal/pkg/hashx"
	"github.com/hanselacn/banking-transaction/internal/pkg/response"
	"github.com/pkg/errors"
)

type UsersHandler struct {
	business business.Business
}

func NewUsersHandler(db *sql.DB) UsersHandler {
	return UsersHandler{business: business.NewBusiness(db)}
}

func (h *UsersHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.users.create_user"
		payload   entity.CreateUserInput
	)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "request body malformed", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	payload.Password, err = hashx.HashPassword(payload.Password)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "hashing password error", nil, err, http.StatusInternalServerError)
		return
	}
	user, err := h.business.UserBusiness.CreateUser(ctx, payload)
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
		response.JsonResponse(w, "create new user error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success create user", user, nil, http.StatusCreated)
}
