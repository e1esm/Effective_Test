package aggregator

import (
	"context"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres"
	"github.com/google/uuid"
)

type Aggregator interface {
	Save(context.Context, *users.ExtendedUser) (uuid.UUID, error)
	Delete(context.Context, uuid.UUID) (uuid.UUID, error)
	Update(context.Context, *users.ExtendedUser) (uuid.UUID, error)
}

type UserService struct {
	repo postgres.Repository
}

func NewUserService(repository postgres.Repository) *UserService {
	return &UserService{
		repo: repository,
	}
}

func (us *UserService) Save(ctx context.Context, user *users.ExtendedUser) (uuid.UUID, error) {
	return us.repo.Save(ctx, *user)
}

func (us *UserService) Delete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	return us.repo.Delete(ctx, id)
}

func (us *UserService) Update(ctx context.Context, user *users.ExtendedUser) (uuid.UUID, error) {
	return us.repo.Update(ctx, *user)
}
