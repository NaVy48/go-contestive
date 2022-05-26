package contest

import (
	"contestive/api/middleware"
	"contestive/entity"
	"context"
)

//Repository interface
type Repository interface {
	Create(ctx context.Context, contest *entity.Contest) error
	GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Contest, error)
	ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Contest, int, error)
	ListAllForUsers(ctx context.Context, options entity.ListOptions, userid entity.ID) ([]entity.Contest, int, error)
	Update(ctx context.Context, contest *entity.Contest) error
	Delete(ctx context.Context, id entity.ID) error
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

func (s Service) GetByID(ctx context.Context, id entity.ID) (entity.Contest, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.GetByID(ctx, id, 0)
	} else {
		return s.Repository.GetByID(ctx, id, uc.UserID)
	}
}

func (s Service) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Contest, int, error) {
	uc := middleware.ClaimsFromContext(ctx)
	if uc.Admin {
		return s.Repository.ListAll(ctx, options)
	} else {
		return s.Repository.ListAllForUsers(ctx, options, uc.UserID)
	}
}
