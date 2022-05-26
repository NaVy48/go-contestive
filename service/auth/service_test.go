package auth

import (
	"contestive/entity"
	"contestive/service/password"
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

var hashedPassword = "$2a$10$OwjCOmEq7jd5Rc5sg3bOFOwAQmx8/xXx/Mbyt3.2jem.rxs9Imo16"

type userServiceStub struct{}

func (us userServiceStub) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	switch username {
	case "user1":
		return entity.User{Username: "user1", PasswordHash: hashedPassword}, nil
	default:
		return entity.User{}, entity.ErrNotFound(fmt.Errorf("repository error"))
	}
}

type jwtServiceStub struct{}

func (s jwtServiceStub) Encode(entity.UserClaims) (entity.Token, error) {
	return "jwt_token", nil
}
func (s jwtServiceStub) Verify(token entity.Token) (uc entity.UserClaims, err error) {
	if token == "jwt_token" {
		uc.UserID = 1
		uc.Admin = false
		uc.Username = "user1"
		return
	}

	err = entity.ErrInvalidToken(nil)
	return
}

func newServiceForTest() *Service {
	return NewService(userServiceStub{}, jwtServiceStub{}, password.NewService())
}

func Test_LogIn(t *testing.T) {
	is := is.New(t)

	s := newServiceForTest()
	token, err := s.LogIn(context.Background(), "user1", "password")
	is.NoErr(err)
	is.True(len(token) > 0)
}

func Test_LogIn_InvalidUser(t *testing.T) {
	is := is.New(t)

	s := newServiceForTest()
	token, err := s.LogIn(context.Background(), "user2", "password")
	is.Equal(err, entity.ErrBadCredentials)
	is.True(len(token) == 0)
}

func Test_LogIn_InvalidPasword(t *testing.T) {
	is := is.New(t)

	s := newServiceForTest()
	token, err := s.LogIn(context.Background(), "user1", "incorrect")
	is.Equal(err, entity.ErrBadCredentials)
	is.True(len(token) == 0)
}
