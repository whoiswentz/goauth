package auth

import (
	"encoding/json"
	"net/http"

	"github.com/whoiswentz/goauth/database"
	"github.com/whoiswentz/goauth/helpers"
	"github.com/whoiswentz/goauth/users"
)

type authHandler struct {
	db *database.Database
}

func NewAuthHandler(db *database.Database) *authHandler {
	return &authHandler{db: db}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := users.ByEmail(h.db, login.Email)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := user.ComparePasswords(login.Password); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	token, refresh, err := helpers.GenerateAllTokens(user.Email, user.Name, user.Id)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	loginRespBytes, err := json.Marshal(LoginResponse{
		Token:        token,
		RefreshToken: refresh,
	})
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(loginRespBytes)
}
