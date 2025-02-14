package repository

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/adapter/repository/storage"
	"github.com/murilo05/JobScheduler/internal/core/domain"
	"go.uber.org/zap"
)

type Repository struct {
	db     storage.Postgres
	aws    storage.AWS
	logger *zap.SugaredLogger
}

func NewRepository(db storage.Postgres, aws storage.AWS, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db,
		aws,
		logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return r.db.Save(ctx, user)
}

func (r *Repository) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	return r.db.List(ctx, skip, limit)
}

func (r *Repository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	return r.db.Get(ctx, id)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.db.GetByEmail(ctx, email)
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return r.db.Update(ctx, user)
}

func (r *Repository) DeleteUser(ctx context.Context, id uint64) error {
	return r.db.Delete(ctx, id)
}

func (r *Repository) SendEmailValidationToSQS(ctx context.Context, email, token string) error {
	return r.aws.SendEmailValidationToSQS(ctx, email, token)
}
