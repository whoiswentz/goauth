package auth

import (
	"encoding/json"
	"net/http"

	"github.com/whoiswentz/goauth/database"
	"github.com/whoiswentz/goauth/users"
)

type authHandler struct {
	db *database.Database
}

func NewAuthHandler(db *database.Database) *authHandler {
	return &authHandler{db: db}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var login Login
	if err := decoder.Decode(&login); err != nil {
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
