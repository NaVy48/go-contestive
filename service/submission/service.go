package submission

import (
	"contestive/api/middleware"
	"contestive/entity"
	"context"
)

//Repository interface
type Repository interface {
	Create(ctx context.Context, submission *entity.Submission) error
	GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Submission, error)
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Submission, int, error)
	ListAllForUsers(ctx context.Context, options entity.ListOptions, userid entity.ID) ([]entity.Submission, int, error)
	UpdateSubmission(ctx context.Context, submission *entity.Submission) error
	GetPendingSubmission(ctx context.Context) (entity.Submission, error)
}

type ProblemService interface {
	GetByID(ctx context.Context, id entity.ID) (entity.Problem, error)
}

type ContestService interface {
	GetByID(ctx context.Context, id entity.ID) (entity.Contest, error)
}

type Notifyer interface {
	Notify()
}

// Service  interface
type Service struct {
	Repository
	jm             Notifyer
	problemService ProblemService
	contestService ContestService
}

// NewService create new use case
func NewService(r Repository, jm Notifyer, problemService ProblemService, contestService ContestService) *Service {
	return &Service{
		r, jm, problemService, contestService,
	}
}

func (s Service) Create(ctx context.Context, submission *entity.Submission) error {
	_, err := s.contestService.GetByID(ctx, submission.ContestID)
	if err != nil {
		return err
	}
	prob, err := s.problemService.GetByID(ctx, submission.ProblemID)
	if err != nil {
		return err
	}
	rev, ok := prob.ActiveRevision()
	if !ok {
		return entity.ErrNotFound(nil)
	}

	uc := middleware.ClaimsFromContext(ctx)
	submission.AuthorID = uc.UserID
	submission.ProblemRevID = rev.ID
	submission.Status = entity.SubmissionStatusPending

	err = s.Repository.Create(ctx, submission)
	if err != nil {
		return err
	}

	s.jm.Notify()
	return nil
}

func (s Service) GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Submission, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.GetByID(ctx, id, 0)
	} else {
		return s.Repository.GetByID(ctx, id, uc.UserID)
	}
}

func (s Service) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Submission, int, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.ListAll(ctx, options)
	} else {
		return s.Repository.ListAllForUsers(ctx, options, uc.UserID)
	}
}
