package problem

import (
	"contestive/api/middleware"
	"contestive/entity"
	"context"
)

//Repository interface
type Repository interface {

	// // ListAll gets users according to options from all active user list
	ListAll(ctx context.Context, options entity.ListOptions, userID entity.ID) ([]entity.Problem, int, error)

	// GetByID gets user by id
	GetByID(ctx context.Context, id, userID entity.ID) (entity.Problem, error)

	// Create creates new user and updates fields that were changed
	// If problem has revisions those are saved too
	Create(ctx context.Context, e *entity.Problem) error

	// Adds new revision to the existing problem. Marks other revisions as outdated.
	AddRevision(ctx context.Context, newRevision *entity.ProblemRevision) error

	StatmentHtmlByProblemId(ctx context.Context, problemID entity.ID) (string, error)
}

// Service  interface
type Service struct {
	Repository
}

// NewService create new use case
func NewService(r Repository) *Service {
	return &Service{
		r,
	}
}

// Create creates new problem from package zip archive
func (s Service) Create(ctx context.Context, archive []byte) (entity.Problem, error) {
	pa, err := ProblemArchive(archive)
	if err != nil {
		return entity.Problem{}, err
	}

	p := pa.Problem
	uc := middleware.ClaimsFromContext(ctx)

	existingProblems, _, err := s.Repository.ListAll(ctx, entity.ListOptions{
		Range:  entity.RangeParam{From: 0, To: 1},
		Filter: []entity.FilterParam{{Field: "externalurl", FilterValues: []interface{}{p.Url}}},
	}, 0)
	if err != nil {
		return entity.Problem{}, err
	}
	if len(existingProblems) > 0 {
		return entity.Problem{}, entity.ErrCustomWrapper("problem already exist", entity.ErrDataConflict(err))
	}

	statementHtml, err := pa.StatementHtml()
	if err != nil {
		return entity.Problem{}, entity.ErrCustomWrapper("missing html statements", err)
	}

	statementPdf, err := pa.StatementPdf()
	if err != nil {
		return entity.Problem{}, entity.ErrCustomWrapper("missing pdf statements", err)
	}

	newProblem := entity.Problem{
		Entity:      entity.Entity{},
		AuthorID:    uc.UserID,
		Name:        p.ShortName,
		ExternalURL: p.Url,
		Revisions: []entity.ProblemRevision{{
			AuthorID:       uc.UserID,
			Revision:       p.Revision,
			Title:          p.Title(),
			MemoryLimit:    p.Testset.MemoryLimit,
			TimeLimit:      p.Testset.TimeLimit,
			StatementHtml:  string(statementHtml),
			StatementPdf:   statementPdf,
			PackageArchive: archive,
			Outdated:       false,
		}},
	}

	err = s.Repository.Create(ctx, &newProblem)
	if err != nil {
		return entity.Problem{}, err
	}

	return newProblem, nil
}

// Update updates user
func (s Service) Update(ctx context.Context, problemID entity.ID, archive []byte) (entity.Problem, error) {
	p, err := s.GetByID(ctx, problemID)
	if err != nil {
		return entity.Problem{}, entity.ErrCustomWrapper("problem retrieving failed", err)
	}

	pa, err := ProblemArchive(archive)
	if err != nil {
		return entity.Problem{}, err
	}

	if pa.Problem.Url != p.ExternalURL {
		return entity.Problem{}, entity.ErrPackageNotCompatible(err)
	}

	if rev, ok := p.ActiveRevision(); ok && pa.Problem.Revision <= rev.Revision {
		return entity.Problem{}, entity.ErrPackageNotCompatible(err)
	}

	uc := middleware.ClaimsFromContext(ctx)

	statementHtml, err := pa.StatementHtml()
	if err != nil {
		return entity.Problem{}, entity.ErrCustomWrapper("missing html statements", err)
	}

	statementPdf, err := pa.StatementPdf()
	if err != nil {
		return entity.Problem{}, entity.ErrCustomWrapper("missing pdf statements", err)
	}

	newRevision := entity.ProblemRevision{
		AuthorID:       uc.UserID,
		Revision:       pa.Problem.Revision,
		Title:          pa.Problem.Title(),
		MemoryLimit:    pa.Problem.Testset.MemoryLimit,
		TimeLimit:      pa.Problem.Testset.TimeLimit,
		StatementHtml:  string(statementHtml),
		StatementPdf:   statementPdf,
		PackageArchive: archive,
		Outdated:       false,
	}

	err = s.Repository.AddRevision(ctx, &newRevision)
	if err != nil {
		return entity.Problem{}, err
	}

	p.Revisions = []entity.ProblemRevision{newRevision}
	return p, nil
}

func (s Service) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Problem, int, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.ListAll(ctx, options, 0)
	} else {
		return s.Repository.ListAll(ctx, options, uc.UserID)
	}

}
func (s Service) GetByID(ctx context.Context, id entity.ID) (entity.Problem, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.GetByID(ctx, id, 0)
	} else {
		return s.Repository.GetByID(ctx, id, uc.UserID)
	}

}
