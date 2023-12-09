package domain

import (
	"golang.org/x/crypto/bcrypt"
)

// User example
type User struct {
	ID       int64
	Username string
	Password string
}

func NewUser(username, password string) User {
	user := User{
		Username: username,
		Password: password,
	}

	return user
}

func (user *User) HashPassword() error {
	newPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(newPassword)

	return nil
}

func (user *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
