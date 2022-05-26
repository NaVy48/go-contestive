package handler

import (
	"contestive/api/apierror"
	"contestive/entity"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type SubmissionService interface {
	Create(ctx context.Context, contest *entity.Submission) error
	GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Submission, error)
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Submission, int, error)
}

type submissionHandler struct {
	ParseResponser
	us     SubmissionService
	router *chi.Mux
}

func NewSubmissionHandler(pr ParseResponser, us SubmissionService) *chi.Mux {
	h := submissionHandler{pr, us, chi.NewRouter()}

	h.router.Get("/", h.list)
	h.router.Post("/", h.create)
	h.router.Get("/{submissionID}", h.getSingle)

	return h.router
}

type submissionRequest struct {
	ProblemID  int64  `json:"problemId"`
	ContestID  string `json:"contestId"`
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

type submissionResponse struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	ProblemID    int64     `json:"problemId"`
	ProblemRevID int64     `json:"problemRevId"`
	ContestID    int64     `json:"contestId"`
	AuthorID     int64     `json:"authorId"`
	Language     string    `json:"language"`
	SourceCode   string    `json:"sourceCode"`
	Status       string    `json:"status"`
	Result       string    `json:"result"`
	Details      string    `json:"details"`
}

func newSubmissionDetailDto(u entity.Submission) submissionResponse {
	return submissionResponse{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		ProblemID:    u.ProblemID,
		ProblemRevID: u.ProblemRevID,
		ContestID:    u.ContestID,
		AuthorID:     u.AuthorID,
		Language:     u.Language,
		SourceCode:   u.SourceCode,
		Status:       string(u.Status),
		Result:       u.Result,
		Details:      u.Details,
	}
}

func (h submissionHandler) list(w http.ResponseWriter, r *http.Request) {
	listop := h.ParseQuery(r)
	listop.Range.LimitAndDefault(1000, 20)

	submissions, total, err := h.us.ListAll(r.Context(), listop)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}
	h.SetTotalCount(w, r, total)

	res := make([]submissionResponse, len(submissions))
	for i, u := range submissions {
		res[i] = newSubmissionDetailDto(u)
	}

	h.ResponseJSON(w, r, res)
}

func (h submissionHandler) getSingle(w http.ResponseWriter, r *http.Request) {
	submissionIdParam := chi.URLParam(r, "submissionID")
	id, err := strconv.Atoi(submissionIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid submission id (%s) in get single submission request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	submission, err := h.us.GetByID(r.Context(), int64(id), 0)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newSubmissionDetailDto(submission))
}

func (h submissionHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload submissionRequest
	err := h.ParseJSON(r, &payload)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	conID, _ := strconv.ParseInt(payload.ContestID, 10, 64)

	submission := entity.Submission{
		ProblemID:  payload.ProblemID,
		ContestID:  conID,
		Language:   payload.Language,
		SourceCode: payload.SourceCode,
	}
	err = h.us.Create(r.Context(), &submission)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newSubmissionDetailDto(submission))
}
