package entity

import "fmt"

var ErrInvalidUserClaims = fmt.Errorf("invalid user claims")

type UserClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	UserID   int64  `json:"uid"`
}
