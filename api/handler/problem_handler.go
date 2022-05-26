package handler

import (
	"bytes"
	"contestive/api/apierror"
	"contestive/api/middleware"
	"contestive/entity"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ProblemService interface {
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Problem, int, error)
	GetByID(ctx context.Context, id entity.ID) (entity.Problem, error)
	Create(ctx context.Context, archive []byte) (entity.Problem, error)
	Update(ctx context.Context, problemID entity.ID, archive []byte) (entity.Problem, error)
	StatmentHtmlByProblemId(ctx context.Context, problemID entity.ID) (string, error)
}

type problemHandler struct {
	ParseResponser
	ps     ProblemService
	router *chi.Mux
}

func NewProblemHandler(pr ParseResponser, us ProblemService) *chi.Mux {
	h := problemHandler{pr, us, chi.NewRouter()}

	h.router.Get("/", h.list)
	h.router.With(middleware.Admin()).Post("/", h.create)
	h.router.Get("/{problemID}", h.getSingle)
	h.router.With(middleware.Admin()).Put("/{problemID}", h.update)
	// h.router.Put("/{problemID}/statement", h.getStatement)

	return h.router
}

type ProblemResponse struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AuthorID    int64     `json:"authorId"`
	ExternalURL string    `json:"externalURL"`
	Name        string    `json:"name"`
	ProblemID   int64     `json:"problemId"`
	Revision    int       `json:"revision"`
	Title       string    `json:"title"`
	MemoryLimit int       `json:"memoryLimit"`
	TimeLimit   int       `json:"timeLimit"`
}

func newProblemResponse(p entity.Problem) ProblemResponse {
	rev, _ := p.ActiveRevision()
	return ProblemResponse{
		ID:          p.ID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		AuthorID:    p.AuthorID,
		ExternalURL: p.ExternalURL,
		Name:        p.Name,
		Revision:    rev.Revision,
		Title:       rev.Title,
		MemoryLimit: rev.MemoryLimit,
		TimeLimit:   rev.TimeLimit,
	}
}

type problemDetailResponse struct {
	ProblemResponse
	StatmentHtml string `json:"statmentHtml"`
}

func newProblemDetailResponse(p entity.Problem) problemDetailResponse {
	rev, _ := p.ActiveRevision()
	return problemDetailResponse{
		newProblemResponse(p),
		rev.StatementHtml,
	}
}

func (h problemHandler) list(w http.ResponseWriter, r *http.Request) {
	listop := h.ParseQuery(r)
	listop.Range.LimitAndDefault(1000, 20)
	h.SetTotalCount(w, r, 0)

	problems, total, err := h.ps.ListAll(r.Context(), listop)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}
	h.SetTotalCount(w, r, total)

	res := make([]ProblemResponse, len(problems))
	for i, u := range problems {
		res[i] = newProblemResponse(u)
	}

	h.ResponseJSON(w, r, res)
}

func (h problemHandler) getSingle(w http.ResponseWriter, r *http.Request) {
	problemIdParam := chi.URLParam(r, "problemID")
	id, err := strconv.Atoi(problemIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid problem id (%s) in get single problem request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	problem, err := h.ps.GetByID(r.Context(), int64(id))
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newProblemDetailResponse(problem))
}

func (h problemHandler) parseFileFromMultipart(r *http.Request, formName string) ([]byte, error) {
	mReader, err := r.MultipartReader()
	if err != nil {
		return nil, apierror.NewHttpErrorWrap(err, "problem create: multipart form read faild", http.StatusBadRequest)
	}

	for {
		part, err := mReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, apierror.NewHttpErrorWrap(err, "problem create: multipart unexpected part", http.StatusBadRequest)
		}

		if part.FormName() == formName {
			const initialBufSize = 5_000_000 // usually should need about 5MB
			const maxFileSize = 20_000_000   // Limit to 20MB

			buf := &bytes.Buffer{}
			buf.Grow(initialBufSize)
			io.Copy(buf, io.LimitReader(part, maxFileSize))
			part.Close()

			archive := buf.Bytes()
			if len(archive) == maxFileSize {
				return nil, apierror.NewHttpErrorWrap(err, "problem create: too big file", http.StatusBadRequest)
			} else {
				return archive, nil
			}
		} else {
			part.Close()
		}
	}

	return nil, apierror.NewHttpErrorWrap(err, fmt.Sprintf("problem create: missing %s form field", formName), http.StatusBadRequest)
}

func (h problemHandler) create(w http.ResponseWriter, r *http.Request) {
	archive, err := h.parseFileFromMultipart(r, "package")
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	problem, err := h.ps.Create(r.Context(), archive)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, problem)
}

func (h problemHandler) update(w http.ResponseWriter, r *http.Request) {
	problemIdParam := chi.URLParam(r, "problemID")
	id, err := strconv.Atoi(problemIdParam)
	if err != nil {
		apiErr := apierror.NewHttpErrorWrap(err, "invalid problem id (%s) in get single problem request", http.StatusBadRequest)
		h.ResponseJSON(w, r, apiErr)
		return
	}

	archive, err := h.parseFileFromMultipart(r, "package")
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	problem, err := h.ps.Update(r.Context(), int64(id), archive)
	if err != nil {
		h.ResponseJSON(w, r, err)
		return
	}

	h.ResponseJSON(w, r, newProblemResponse(problem))
}

// func (h problemHandler) getStatement(w http.ResponseWriter, r *http.Request) {
// 	problemIdParam := chi.URLParam(r, "problemID")
// 	id, err := strconv.Atoi(problemIdParam)
// 	if err != nil {
// 		apiErr := apierror.NewHttpErrorWrap(err, "invalid problem id (%s) in get single problem request", http.StatusBadRequest)
// 		h.ResponseJSON(w, r, apiErr)
// 		return
// 	}

// 	statement, err := h.ps.StatmentHtmlByProblemId(r.Context(), int64(id))
// 	if err != nil {
// 		h.ResponseJSON(w, r, err)
// 		return
// 	}

// 	io.Copy(w, strings.NewReader(statement))
// }
