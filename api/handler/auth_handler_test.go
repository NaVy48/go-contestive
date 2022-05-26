package handler

import (
	"contestive/api/payload"
	"contestive/entity"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func makePayload(username, password string) string {
	return fmt.Sprintf(`{"username": "%s","password": "%s"}`, username, password)
}

type AuthServiceStub struct {
	token string
}

func (s *AuthServiceStub) LogIn(ctx context.Context, username, password string) (entity.Token, error) {
	if username != "username" || password != "password" {
		return "", entity.ErrBadCredentials(fmt.Errorf("repository error"))
	}
	return s.token, nil
}

const expextedToken = "this is expected token"

var handler = authHandler{payload.NewJSONHandler(log.New(io.Discard, "", 0)), &AuthServiceStub{expextedToken}}

func TestLogin_BadRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		// TODO: Add test cases.
		{"payload null", "null"},
		{"payload wrong user", makePayload("user1", "pass123")},
		{"payload wrong password", makePayload("username", "pass123")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(tt.payload))
			handler.Login().ServeHTTP(w, r)
			is.Equal(w.Code, http.StatusBadRequest) // Status should be 400

			response := struct {
				Error string `json:"error"`
			}{}

			err := json.NewDecoder(w.Result().Body).Decode(&response)
			is.NoErr(err)
			is.True(len(response.Error) > 10)

		})
	}
}

func TestLogin(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(makePayload("username", "password")))
	handler.Login().ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK) // Status should be 200

	response := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(w.Result().Body).Decode(&response)
	is.NoErr(err)
	is.Equal(expextedToken, response.Token)
}
