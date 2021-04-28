package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/whoiswentz/goauth/sqls"
)

type Database struct {
	Db *sql.DB
}

func Open() (*Database, error) {
	db, err := sql.Open("sqlite3", "database.db")
	return &Database{Db: db}, err
}

func (d Database) RunMigrations() {
	_, err := d.Db.Exec(sqls.POSTS_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}

	_, err = d.Db.Exec(sqls.USERS_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
}
