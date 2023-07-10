package models

import (
	"errors"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	Root      bool      `json:"root,omitempty"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// HashPassword hashes a User Password
func (u *User) HashPassword() error {
	if len(u.Password) != 0 {
		password := []byte(u.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
		return nil
	}
	return errors.New("user password is missing")
}
