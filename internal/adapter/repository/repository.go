package repository

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/adapter/repository/storage"
	"github.com/murilo05/JobScheduler/internal/core/domain"
	"go.uber.org/zap"
)

type Repository struct {
	storage storage.UserStorage
	logger  *zap.SugaredLogger
}

func NewRepository(storage storage.UserStorage, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		storage,
		logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return r.storage.Save(ctx, user)
}

func (r *Repository) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	return r.storage.List(ctx, skip, limit)
}

func (r *Repository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	return r.storage.Get(ctx, id)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.storage.GetByEmail(ctx, email)
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return r.storage.Update(ctx, user)
}

func (r *Repository) DeleteUser(ctx context.Context, id uint64) error {
	return r.storage.Delete(ctx, id)
}

func (r *Repository) SendEmailValidationToSQS(ctx context.Context, email, token string) error {
	return nil
}
