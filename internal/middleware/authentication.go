package middleware

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hanselacn/banking-transaction/internal/pkg/hashx"
	"github.com/hanselacn/banking-transaction/internal/pkg/response"
	"github.com/hanselacn/banking-transaction/repo"
)

type Middleware interface {
	AuthenticationMiddleware(next http.Handler, roles ...string) http.Handler
}

type middleware struct {
	repo repo.Repo
}

func NewMiddleware(db *sql.DB) Middleware {
	return &middleware{
		repo: repo.NewRepositories(db),
	}
}

func (m *middleware) AuthenticationMiddleware(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx       = r.Context()
			eventName = "middleware.authentication"
		)
		log.Println("Authentication check")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			response.JsonResponse(w, "Unauthorized", nil, "unauthorized", http.StatusUnauthorized)
			return
		}
		authValue := strings.SplitN(authHeader, " ", 2)
		if len(authValue) != 2 || authValue[0] != "Basic" {
			response.JsonResponse(w, "Unauthorized", nil, "unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(authValue[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			response.JsonResponse(w, "Unauthorized", nil, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.repo.Users.FindByUserName(ctx, pair[0])
		if err != nil {
			log.Println(eventName, err)
			response.JsonResponse(w, "Unauthorized", nil, "unauthorized", http.StatusUnauthorized)
			return
		}

		var foundRole bool
		for i := range roles {
			if roles[i] == user.Role {
				foundRole = true
			}
		}

		if !foundRole {
			response.JsonResponse(w, "Forbidden", nil, "You Don't Have Access to this Feature! Please contact your Super Admin.", http.StatusForbidden)
			return
		}

		auth, err := m.repo.Authorization.FindByUserID(ctx, user.ID)
		if err != nil {
			log.Println(eventName, err)
			response.JsonResponse(w, "Unauthorized", nil, "User Not Found!", http.StatusUnauthorized)
			return
		}

		match := hashx.CheckPasswordHash(pair[1], auth.Password)
		if !match {
			response.JsonResponse(w, "Forbidden", nil, "Wrong Password", http.StatusForbidden)
			return
		}
		r = r.WithContext(context.WithValue(ctx, CtxValueUserName, user.Username))
		r = r.WithContext(context.WithValue(r.Context(), CtxValueRole, user.Role))
		log.Println("User Authenticated!")
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
