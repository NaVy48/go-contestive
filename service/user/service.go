package user

import (
	"contestive/entity"
	"context"
	"time"
)

//Repository interface
type Repository interface {
	// GetByUsername gets user by username
	GetByUsername(ctx context.Context, username string) (entity.User, error)

	// ListAll gets users according to options from all active user list
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.User, int, error)

	// GetByID gets user by id
	GetByID(ctx context.Context, id entity.ID) (entity.User, error)

	// Create creates new user and updates fields that were changed
	Create(ctx context.Context, e *entity.User) error

	// Update updates user and updates fields hat were changed
	Update(ctx context.Context, e *entity.User) error

	// Update updates user and updates fields hat were changed
	Delete(ctx context.Context, id entity.ID) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
}

// Service  interface
type Service struct {
	Repository
	passwordHasher PasswordHasher
}

// NewService create new use case
func NewService(r Repository, ph PasswordHasher) *Service {
	return &Service{
		r,
		ph,
	}
}

// Create creates new user
func (s Service) Create(ctx context.Context, username, firstName, lastName, password string, isAdmin bool) (entity.User, error) {
	var user entity.User
	hashedPwd, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		return user, entity.ErrCustomWrapper("create failed while hashing password", err)
	}

	user, err = entity.NewUser(username, firstName, lastName, hashedPwd, isAdmin)
	if err != nil {
		return user, err
	}

	err = s.Repository.Create(ctx, &user)
	if err != nil {
		return user, entity.ErrCustomWrapper("user creation failed", err)
	}

	return user, nil
}

// Update updates user
func (s Service) Update(ctx context.Context, id entity.ID, username, firstName, lastName, password string, isAdmin bool) (entity.User, error) {
	var user entity.User

	user, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return user, entity.ErrCustomWrapper("user update failed", err)
	}

	user.Username = username
	user.FirstName = firstName
	user.LastName = lastName
	user.Admin = isAdmin
	user.UpdatedAt = time.Now()

	if password != "" {
		hashedPwd, err := s.passwordHasher.HashPassword(password)
		if err != nil {
			return user, entity.ErrCustomWrapper("create failed while hashing password", err)
		}
		user.PasswordHash = hashedPwd
	}

	user.Validate()
	if err != nil {
		return user, entity.ErrCustomWrapper("update validation failed", err)
	}

	err = s.Repository.Update(ctx, &user)
	if err != nil {
		return user, entity.ErrCustomWrapper("user update failed", err)
	}

	return user, nil
}
