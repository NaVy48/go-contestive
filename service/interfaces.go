package service

import "contestive/entity"

//UserService interface
type UserService interface {
	GetUserByUsername(username string) (entity.User, error)
	// GetUser(id int64) (entity.User, error)
	// SearchUsers(query string) ([]entity.User, error)
	// ListUsers() ([]entity.User, error)
	// CreateUser(username, email, password, firstName, lastName string) error
	// UpdateUser(e entity.User) error
	// DeleteUser(id int64) error
}
