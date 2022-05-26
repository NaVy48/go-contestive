package handler

import (
	"contestive/api/apierror"
	"contestive/api/middleware"
	"contestive/entity"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	// ListAll gets users according to options from all active user list
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.User, int, error)
	GetByID(ctx context.Context, id entity.ID) (entity.User, error)
	Create(ctx context.Context, username, firstName, lastName, password string, isAdmin bool) (entity.User, error)
	Update(ctx context.Context, id entity.ID, username, firstName, lastName, password string, isAdmin bool) (entity.User, error)
	Delete(ctx context.Context, id entity.ID) error
}

type userHandler struct {
	ParseResponser
	us     UserService
	router *chi.Mux
}

func NewUserHandler(pr ParseResponser, us UserService) *chi.Mux {
	h := userHandler{pr, us, chi.NewRouter()}

	h.router.Use(middleware.Admin())
	h.router.Get("/", h.list)
	h.router.Post("/", h.create)
	h.router.Get("/{userID}", h.getSingle)
	h.router.Put("/{userID}", h.update)
	h.router.Delete("/{userID}", h.delete)

	return h.router
}

type userRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Admin     bool   `json:"admin"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Admin     bool      `json:"admin"`
}

func newUserDetailDto(u entity.User) userResponse {
	return userResponse{
		u.ID,
		u.Username,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.UpdatedAt,
		u.Admin,
	}
}

func (h userHandler) list(w http.ResponseWriter, r *http.Request) {
	listop := h.ParseQuery(r)
	listop.Range.LimitAndDefault(1000, 20)

	users, total, err := h.us.ListAll(r.Context(), listop)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}
	h.SetTotalCount(w, r, total)

	res := make([]userResponse, len(users))
	for i, u := range users {
		res[i] = newUserDetailDto(u)
	}

	h.ResponseJSON(w, r, res)
}

func (h userHandler) getSingle(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid user id (%s) in get single user request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	user, err := h.us.GetByID(r.Context(), int64(id))
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newUserDetailDto(user))
}

func (h userHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload userRequest
	err := h.ParseJSON(r, &payload)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	user, err := h.us.Create(r.Context(), payload.Username, payload.FirstName, payload.LastName, payload.Password, payload.Admin)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newUserDetailDto(user))
}

func (h userHandler) update(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid user id (%s) in get single user request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	var payload userRequest
	err = h.ParseJSON(r, &payload)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	user, err := h.us.Update(r.Context(), int64(id), payload.Username, payload.FirstName, payload.LastName, payload.Password, payload.Admin)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newUserDetailDto(user))
}

func (h userHandler) delete(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid user id (%s) in get single user request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	err = h.us.Delete(r.Context(), int64(id))
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
