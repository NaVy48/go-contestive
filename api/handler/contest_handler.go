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

type ContestService interface {
	Create(ctx context.Context, contest *entity.Contest) error
	GetByID(ctx context.Context, id entity.ID) (entity.Contest, error)
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Contest, int, error)
	Update(ctx context.Context, contest *entity.Contest) error
	Delete(ctx context.Context, id entity.ID) error
}

type constestHandler struct {
	ParseResponser
	us     ContestService
	router *chi.Mux
}

func NewContestHandler(pr ParseResponser, us ContestService) *chi.Mux {
	h := constestHandler{pr, us, chi.NewRouter()}

	h.router.Get("/", h.list)
	h.router.With(middleware.Admin()).Post("/", h.create)
	h.router.Get("/{constestID}", h.getSingle)
	h.router.With(middleware.Admin()).Put("/{constestID}", h.update)
	h.router.With(middleware.Admin()).Delete("/{constestID}", h.delete)

	return h.router
}

type constestRequest struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Problems  []int64   `json:"problems"`
	Users     []int64   `json:"users"`
}

type constestResponse struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ProblemID   int64     `json:"problemId"`
	AuthorID    int64     `json:"authorId"`
	ContestName string    `json:"title"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Problems    []int64   `json:"problems"`
	Users       []int64   `json:"users"`
}

func newContestDetailDto(u entity.Contest) constestResponse {
	return constestResponse{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		ProblemID:   u.ProblemID,
		AuthorID:    u.AuthorID,
		ContestName: u.ContestName,
		StartTime:   u.StartTime,
		EndTime:     u.EndTime,
		Problems:    u.Problems,
		Users:       u.Users,
	}
}

func (h constestHandler) list(w http.ResponseWriter, r *http.Request) {
	listop := h.ParseQuery(r)
	listop.Range.LimitAndDefault(1000, 20)

	constests, total, err := h.us.ListAll(r.Context(), listop)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}
	h.SetTotalCount(w, r, total)

	res := make([]constestResponse, len(constests))
	for i, u := range constests {
		res[i] = newContestDetailDto(u)
	}

	h.ResponseJSON(w, r, res)
}

func (h constestHandler) getSingle(w http.ResponseWriter, r *http.Request) {
	constestIdParam := chi.URLParam(r, "constestID")
	id, err := strconv.Atoi(constestIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid constest id (%s) in get single constest request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	constest, err := h.us.GetByID(r.Context(), int64(id))
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newContestDetailDto(constest))
}

func (h constestHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload constestRequest
	err := h.ParseJSON(r, &payload)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	uc := middleware.ClaimsFromContext(r.Context())

	constest := entity.Contest{
		ContestName: payload.Title,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		AuthorID:    uc.UserID,
	}
	err = h.us.Create(r.Context(), &constest)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newContestDetailDto(constest))
}

func (h constestHandler) update(w http.ResponseWriter, r *http.Request) {
	constestIdParam := chi.URLParam(r, "constestID")
	id, err := strconv.Atoi(constestIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid constest id (%s) in get single constest request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	var payload constestRequest
	err = h.ParseJSON(r, &payload)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	constest := entity.Contest{
		Entity:      entity.Entity{ID: int64(id)},
		ContestName: payload.Title,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Problems:    payload.Problems,
		Users:       payload.Users,
	}
	err = h.us.Update(r.Context(), &constest)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newContestDetailDto(constest))
}

func (h constestHandler) delete(w http.ResponseWriter, r *http.Request) {
	constestIdParam := chi.URLParam(r, "constestID")
	id, err := strconv.Atoi(constestIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid constest id (%s) in get single constest request", http.StatusBadRequest)
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
