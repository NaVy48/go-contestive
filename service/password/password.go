package password

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
}

func NewService() PasswordService {
	return passwordService{}
}

func (passwordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (passwordService) VerifyPassword(password, hash string) bool {
	// will be a string so we'll need to convert it to a byte slice
	p, h := []byte(password), []byte(hash)
	err := bcrypt.CompareHashAndPassword(h, p)

	return err == nil
}
