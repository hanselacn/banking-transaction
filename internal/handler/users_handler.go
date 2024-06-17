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
	"github.com/hanselacn/banking-transaction/internal/pkg/hashx"
	"github.com/hanselacn/banking-transaction/internal/pkg/response"
	"github.com/hanselacn/banking-transaction/internal/pkg/rule"
	"github.com/pkg/errors"
)

func MountUserHandler(r *mux.Router, h handler, m middleware.Middleware) {
	r.Handle("/banking-transaction/users/create", m.AuthenticationMiddleware((http.HandlerFunc(h.UsersHandler.CreateUserHandler)), consts.RoleSuperAdmin)).Methods("POST")
	r.Handle("/banking-transaction/users/detail/{user_name}", m.AuthenticationMiddleware((http.HandlerFunc(h.UsersHandler.GetUserDetail)), []string{consts.RoleSuperAdmin, consts.RoleAdmin, consts.RoleCustomer}...)).Methods("GET")
	r.Handle("/banking-transaction/users/update/role", http.HandlerFunc(h.UsersHandler.UpdateRoleByUserName)).Methods("PUT")

	// r.Handle("/banking-transaction/users/create/supadmin", http.HandlerFunc(h.UsersHandler.CreateUserHandler)).Methods("POST")
}

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

	if err := payload.Validate(); err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "create user error", nil, err.Error(), http.StatusUnprocessableEntity)
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

func (h *UsersHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.users.get_detail"
		pathVar   = mux.Vars(r)
		username  = pathVar["user_name"]
	)
	role := ctx.Value(middleware.CtxValueRole)
	switch role {
	case "customer":
		ctxUserName := ctx.Value(middleware.CtxValueUserName)
		if ctxUserName != username {
			response.JsonResponse(w, "Forbidden", nil, "You Have to Access your Own User Information", http.StatusForbidden)
			return
		}
	}

	if err := validation.Validate(username, rule.UserNameRule); err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "get user error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	user, err := h.business.UserBusiness.GetUserDetail(ctx, username)
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
		response.JsonResponse(w, "get user error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success get user", user, nil, http.StatusOK)
}

func (h *UsersHandler) UpdateRoleByUserName(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		eventName = "handler.users.update_role_by_user_name"
		payload   entity.User
	)
	role := ctx.Value(middleware.CtxValueRole)
	switch role {
	case "customer":
		response.JsonResponse(w, "Forbidden", nil, "You Cannot Have Access for this Feature", http.StatusForbidden)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "request body malformed", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := payload.Validate(); err != nil {
		log.Println(eventName, err)
		response.JsonResponse(w, "update user role error", nil, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = h.business.UserBusiness.UpdateRoleByUserName(ctx, payload)
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
		response.JsonResponse(w, "update role got error", nil, err, statusCode)
		return
	}
	response.JsonResponse(w, "success update role", nil, nil, http.StatusOK)
}
