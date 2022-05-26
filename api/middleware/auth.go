package middleware

import (
	"contestive/api/apierror"
	"contestive/entity"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ParseResponser interface {
	ParseQuery(r *http.Request) entity.ListOptions
	ParseJSON(r *http.Request, v interface{}) error
	ResponseJSON(w http.ResponseWriter, r *http.Request, v interface{})
	SetTotalCount(w http.ResponseWriter, r *http.Request, total int)
}

type AuthService interface {
	Validate(token string) (entity.UserClaims, error)
}

type contextKey string

var claimsContextKey contextKey = "userClaims"

// tokenGetter get bearer token form Authorization header
func tokenGetter(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// Auth retrun middlware for handling user authentication
func Auth(p ParseResponser, as AuthService) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := as.Validate(tokenGetter(r))
			if err != nil {
				p.ResponseJSON(w, r, apierror.NewHttpErrorWrap(err, "invalid or missing token", http.StatusUnauthorized))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Admin() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if !claims.Admin {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ClaimsFromContext gets user claims from ctx
func ClaimsFromContext(ctx context.Context) entity.UserClaims {
	fmt.Printf("%[1]T %[1]v\n", ctx.Value(claimsContextKey))

	if claims, ok := ctx.Value(claimsContextKey).(entity.UserClaims); ok {
		return claims
	}

	return entity.UserClaims{}
}
