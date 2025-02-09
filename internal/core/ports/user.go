package ports

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}
