package storage

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/core/domain"
)

type UserStorage interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	List(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	Get(ctx context.Context, id uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id uint64) error
}
