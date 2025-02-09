package service

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/core/domain"
	"github.com/murilo05/JobScheduler/internal/core/ports"
	"github.com/murilo05/JobScheduler/internal/core/util"
	"go.uber.org/zap"
)

type UserService struct {
	repo   ports.UserRepository
	logger *zap.SugaredLogger
}

func NewUserService(repo ports.UserRepository, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		repo,
		logger,
	}
}

func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		us.logger.Error("Failed to hash password: ", err)
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			us.logger.Error("data already exist: ", err)
			return nil, err
		}
		us.logger.Error("failed to create user: ", err)
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	var user *domain.User

	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			us.logger.Error("user not found")
			return nil, err
		}
		us.logger.Error("internal error: %s", err)
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var users []domain.User

	users, err := us.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		us.logger.Error("internal error: %s", err)
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			us.logger.Error("user not found")
			return nil, err
		}
		us.logger.Error("internal error: %s", err)
		return nil, domain.ErrInternal
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = util.HashPassword(user.Password)
		if err != nil {
			us.logger.Error("auth failed: %s", err)
			return nil, domain.ErrInternal
		}
	}

	user.Password = hashedPassword

	_, err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			us.logger.Error("data already exist")
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			us.logger.Error("user not found")
			return err
		}
		return domain.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}
