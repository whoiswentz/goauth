package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/whoiswentz/goauth/app"
	"github.com/whoiswentz/goauth/database"
)

type userHandler struct {
	db *database.Database
}

func NewUserHandler(db *database.Database) *userHandler {
	return &userHandler{db: db}
}

func (h userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		appError := app.NewAppError(w, err, http.StatusInternalServerError)
		w.Write(appError)
		return
	}

	createdUser, err := create(h.db, &user)
	if err != nil {
		appError := app.NewAppError(w, err, http.StatusBadRequest)
		w.Write(appError)
		return
	}

	userBytes, err := json.Marshal(createdUser)
	if err != nil {
		appError := app.NewAppError(w, err, http.StatusInternalServerError)
		w.Write(appError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userBytes)
}

func (h userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := list(h.db)
	if err != nil {
		appError := app.NewAppError(w, err, http.StatusBadRequest)
		w.Write(appError)
		return
	}

	usersByte, err := json.Marshal(users)
	if err != nil {
		appError := app.NewAppError(w, err, http.StatusInternalServerError)
		w.Write(appError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersByte)
}

func (h userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		appError := app.NewAppError(w, err, http.StatusInternalServerError)
		w.Write(appError)
		return
	}

	user, err := byId(h.db, int64(id))
	if err != nil {
		appError := app.NewAppError(w, err, http.StatusBadRequest)
		w.Write(appError)
		return
	}

	if err := delete(h.db, *user); err != nil {
		appError := app.NewAppError(w, err, http.StatusBadRequest)
		w.Write(appError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
