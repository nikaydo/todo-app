package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/nikaydo/final/internal/jwt"
)

var (
	ErrPasswordError = errors.New("wrong password or not set")
)

type Signin struct {
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func (h *Handlers) signin(w http.ResponseWriter, r *http.Request) {
	var pass Signin
	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSONResponse(w, MsgError{Error: ErrPasswordError.Error()}, http.StatusUnauthorized)
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
	if pass.Password == password {
		token, err := jwt.CreateJwt(todo_secret_jwt, todo_hash_password)
		if err != nil {
			writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, Token{Token: token}, http.StatusOK)
		return
	}
	writeJSONResponse(w, MsgError{Error: ErrPasswordError.Error()}, http.StatusUnauthorized)
}
