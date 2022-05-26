// Package jwt provides a service to create and verify
// JWT auth tokens for the bebop web app.
package jwt

import (
	"contestive/entity"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// NewJWTService creates a new JWT service using the given secret (32-byte hex-encoded).
func NewJWTService(secret string, tokenExpiration time.Duration) (*jwtService, error) {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode jwt secret from hex: %s", err)
	}

	if len(secretBytes) < 32 {
		return nil, errors.New("jwt: secret too short")
	}

	return &jwtService{secret: secretBytes, tokenExpiration: tokenExpiration, parser: jwt.NewParser()}, nil
}

// jwtService is service used for handling JWT tokens
type jwtService struct {
	secret          []byte
	tokenExpiration time.Duration
	parser          *jwt.Parser
}

// claims strucure for keeping claims data
type claims struct {
	jwt.StandardClaims
	entity.UserClaims
}

func (c claims) Valid(h *jwt.ValidationHelper) error {
	err := c.StandardClaims.Valid(h)
	if err != nil {
		return err
	}

	if c.UserClaims.UserID == 0 || c.UserClaims.Username == "" {
		return entity.ErrInvalidUserClaims
	}

	return nil
}

// Encode creates a JWT string with provided claims
func (s *jwtService) Encode(userClaims entity.UserClaims) (string, error) {
	now := time.Now()
	claims := claims{
		jwt.StandardClaims{
			IssuedAt:  jwt.At(now),
			ExpiresAt: jwt.At(now.Add(s.tokenExpiration)),
		},
		userClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("jwt: token signing failed: %w", err)
	}

	return tokenString, nil
}

// Verify verifies the JWT string using the given secret key.
// On success it returns the user ID and the time the token was issued.
func (s *jwtService) Verify(tokenString string) (entity.UserClaims, error) {
	c := claims{}
	_, err := s.parser.ParseWithClaims(
		tokenString,
		&c,
		jwt.KnownKeyfunc(jwt.SigningMethodHS256, s.secret),
	)
	if err != nil {
		return entity.UserClaims{}, err
	}
	return c.UserClaims, nil

}

// // Middleware retrun middlware for handling user authentication
// func Middleware(s *JwtService) func(h http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		hfn := func(w http.ResponseWriter, r *http.Request) {
// 			claims, _, err := s.Verify(tokenGetter(r))
// 			if err != nil {
// 				// TODO fix unauthorized error
// 				fmt.Println(err)
// 				w.WriteHeader(401)
// 				return
// 			}

// 			ctx := r.Context()
// 			ctx = context.WithValue(ctx, claimsContextKey, claims)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		}
// 		return http.HandlerFunc(hfn)
// 	}
// }

// // ClaimsFromContext gets user claims from ctx
// func ClaimsFromContext(ctx context.Context) Claims {
// 	fmt.Printf("%[1]T %[1]v\n", ctx.Value(claimsContextKey))

// 	if claims, ok := ctx.Value(claimsContextKey).(Claims); ok {
// 		return claims
// 	}

// 	return Claims{}
// }

// type contextKey string

// // tokenGetter get bearer token form Authorization header
// func tokenGetter(r *http.Request) string {
// 	bearer := r.Header.Get("Authorization")
// 	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
// 		return bearer[7:]
// 	}
// 	return ""
// }
