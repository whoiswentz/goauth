package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/whoiswentz/goauth/auth"
	"github.com/whoiswentz/goauth/infrastructure/cache"

	"github.com/whoiswentz/goauth/database"
	"github.com/whoiswentz/goauth/middlewares"
	"github.com/whoiswentz/goauth/posts"
	"github.com/whoiswentz/goauth/users"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	db.RunMigrations()

	blackListCache := cache.NewCacheWithTTL()
	mux := NewRouter(db, blackListCache)

	port := os.Getenv("PORT")
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func NewRouter(db *database.Database, c *cache.Cache) *mux.Router {
	ph := posts.NewPostsHandler(db)
	uh := users.NewUserHandler(db)
	ah := auth.NewAuthHandler(db)

	mux := mux.NewRouter()

	mux.HandleFunc("/posts/create", ph.CreatePost).Methods(http.MethodPost)
	mux.HandleFunc("/posts/list", ph.ListPosts).Methods(http.MethodGet)

	mux.HandleFunc("/users", uh.CreateUser).Methods(http.MethodPost)
	mux.HandleFunc("/users", uh.ListUsers).Methods(http.MethodGet)
	mux.HandleFunc("/users/{id}", middlewares.Chain(
		uh.DeleteUser,
		middlewares.RequireToken(c),
	)).Methods(http.MethodDelete)

	mux.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)

	return mux
}
