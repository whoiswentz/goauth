package users

import (
	"database/sql"
	"log"

	"github.com/whoiswentz/goauth/database"
)

func create(db *database.Database, u *User) (*User, error) {
	stmt, err := db.Db.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(u.Name, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.Id = id

	return u, nil
}

func byId(db *database.Database, id int64) (*User, error) {
	row := db.Db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func byEmail(db *database.Database, email string) (*User, error) {
	row := db.Db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	var user User
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func list(db *database.Database) ([]User, error) {
	rows, err := db.Db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func delete(db *database.Database, u User) error {
	stmt, err := db.Db.Prepare("DELETE FROM users WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("total of %d rows affected", count)

	return nil
}
