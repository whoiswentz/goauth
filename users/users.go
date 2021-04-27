package users

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) HashPassword() error {
	bytePass := []byte(u.Password)
	hashPass, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashPass)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)
}
