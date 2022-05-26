package handler

import (
	"contestive/api/apierror"
	"contestive/entity"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	LogIn(ctx context.Context, username, password string) (entity.Token, error)
}

type authHandler struct {
	ParseResponser
	as AuthService
}

func NewAuthHandler(pr ParseResponser, as AuthService) *chi.Mux {
	h := authHandler{pr, as}

	router := chi.NewRouter()
	router.Post("/login", h.Login())

	return router
}

func (h authHandler) Login() http.HandlerFunc {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type Res struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var rp Req
		err := json.NewDecoder(r.Body).Decode(&rp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := h.as.LogIn(r.Context(), rp.Username, rp.Password)
		if err != nil {
			if errors.Is(err, entity.ErrBadCredentials(nil)) {
				h.ResponseJSON(w, r, apierror.NewHttpErrorWrap(err, "Invalid username or password", http.StatusBadRequest))
				return
			}

			h.ResponseJSON(w, r, apierror.NewHttpErrorWrap(err, "Internal error during authentication", http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Res{token})
	}
}
