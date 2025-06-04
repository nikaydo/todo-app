package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/nikaydo/final/internal/jwt"
)

var (
	ErrConfiguration = errors.New("server configuration error")
	ErrAuthError     = errors.New("authentication required")
	ErrValidJwt      = errors.New("error validation jwt token")
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := os.Getenv("TODO_PASSWORD")
		if len(secretKey) == 0 {
			writeJSONResponse(w, MsgError{Error: ErrConfiguration.Error()}, http.StatusInternalServerError)
			return
		}
		todo_secret_jwt := os.Getenv("TODO_SECRET_JWT")
		if len(todo_secret_jwt) == 0 {
			writeJSONResponse(w, MsgError{Error: ErrConfiguration.Error()}, http.StatusInternalServerError)
			return
		}
		todo_hash_password := os.Getenv("TODO_HASH_PASSWORD")
		if len(todo_hash_password) == 0 {
			writeJSONResponse(w, MsgError{Error: ErrConfiguration.Error()}, http.StatusInternalServerError)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			writeJSONResponse(w, MsgError{Error: ErrAuthError.Error()}, http.StatusUnauthorized)
			return
		}
		valid, err := jwt.ValidateJwt(cookie.Value, todo_secret_jwt, todo_hash_password)
		if err != nil {
			writeJSONResponse(w, MsgError{Error: ErrValidJwt.Error()}, http.StatusUnauthorized)
			return
		}
		if !valid {
			writeJSONResponse(w, MsgError{Error: ErrAuthError.Error()}, http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}
