package main

import (
	"log"
	"net/http"
	"time"

	"github.com/whoiswentz/goauth/database"
	"github.com/whoiswentz/goauth/posts"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	db.RunMigrations()

	ph := posts.NewPostsHandler(db)

	r := http.NewServeMux()
	r.HandleFunc("/posts/create", ph.CreatePost)
	r.HandleFunc("/posts/list", ph.ListPosts)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
