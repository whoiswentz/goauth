package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/whoiswentz/goauth/auth"
	"github.com/whoiswentz/goauth/database"
	"github.com/whoiswentz/goauth/posts"
	"github.com/whoiswentz/goauth/users"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	db.RunMigrations()

	ph := posts.NewPostsHandler(db)
	uh := users.NewUserHandler(db)
	ah := auth.NewAuthHandler(db)

	mux := mux.NewRouter()
	mux.HandleFunc("/posts/create", ph.CreatePost).Methods(http.MethodPost)
	mux.HandleFunc("/posts/list", ph.ListPosts).Methods(http.MethodGet)

	mux.HandleFunc("/users", uh.CreateUser).Methods(http.MethodPost)
	mux.HandleFunc("/users", uh.ListUsers).Methods(http.MethodGet)
	mux.HandleFunc("/users/:id", uh.DeleteUser).Methods(http.MethodDelete)

	mux.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
