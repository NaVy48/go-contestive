package auth

import (
	"contestive/entity"
	"context"
	"errors"
	"fmt"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (entity.User, error)
}

type JwtService interface {
	Encode(userClaims entity.UserClaims) (string, error)
	Verify(tokenString string) (entity.UserClaims, error)
}

type PasswordVerifyer interface {
	VerifyPassword(password, paswordHash string) bool
}

//Service  interface
type Service struct {
	us  UserService
	jwt JwtService
	pwd PasswordVerifyer
}

//NewService create new use case
func NewService(us UserService, jwt JwtService, pwd PasswordVerifyer) *Service {
	return &Service{us, jwt, pwd}
}

// LogIn checks username and password, and generates a token
func (s *Service) LogIn(ctx context.Context, username, password string) (entity.Token, error) {
	u, err := s.us.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound(nil)) {
			return "", entity.ErrBadCredentials(err)
		}

	}

	ok := s.pwd.VerifyPassword(password, u.PasswordHash)
	if !ok {
		return "", entity.ErrBadCredentials(fmt.Errorf("user %s entered incorrect password", username))
	}

	token, err := s.jwt.Encode(entity.UserClaims{Username: u.Username, Admin: u.Admin, UserID: u.ID})
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate check token validity and returns userID and admin flag
func (s *Service) Validate(token entity.Token) (uc entity.UserClaims, err error) {
	return s.jwt.Verify(token)
}
