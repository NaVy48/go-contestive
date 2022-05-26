package entity

import (
	"fmt"
	"time"
)

var ErrUserMissingFields = fmt.Errorf("user missing some fields")

// User db model
type User struct {
	Entity
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
	Admin        bool
}

//NewUser create a new user
func NewUser(username, firstName, lastName, passwordHash string, isAdmin bool) (u User, err error) {
	u = User{
		Entity{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
		username,
		firstName,
		lastName,
		passwordHash,
		isAdmin,
	}

	err = u.Validate()
	return
}

//Validate validate data
func (u User) Validate() error {
	if u.Username == "" || u.FirstName == "" || u.LastName == "" || u.PasswordHash == "" {
		return ErrUserMissingFields
	}

	return nil
}
