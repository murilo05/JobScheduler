package service

import (
	"context"

	"github.com/murilo05/JobScheduler/internal/core/domain"
	"github.com/murilo05/JobScheduler/internal/core/ports"
	"github.com/murilo05/JobScheduler/internal/core/util"
	"go.uber.org/zap"
)

type AuthService struct {
	repo   ports.Repository
	ts     ports.TokenService
	logger *zap.SugaredLogger
}

func NewAuthService(repo ports.Repository, ts ports.TokenService, logger *zap.SugaredLogger) *AuthService {
	return &AuthService{
		repo,
		ts,
		logger,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			as.logger.Error("invalid credentials")
			return "", domain.ErrInvalidCredentials
		}
		as.logger.Error("auth error: %s", err)
		return "", domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		as.logger.Error("invalid credentials")
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		as.logger.Error("token creation failed: %s", err)
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}
