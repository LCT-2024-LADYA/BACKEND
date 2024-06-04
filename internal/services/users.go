package services

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"BACKEND/pkg/utils"
	"context"
	"github.com/rs/zerolog"
	"time"
)

type userService struct {
	userRepo       repository.Users
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitUserService(
	userRepo repository.Users,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Users {
	return &userService{
		userRepo:       userRepo,
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (u userService) Register(ctx context.Context, user domain.UserCreate) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.dbResponseTime)
	defer cancel()

	createdID, err := u.userRepo.Create(ctx, user)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return 0, err
	}

	u.logger.Info().Msg(log.Normalizer(log.CreateObject, log.User, createdID))

	return createdID, nil
}

func (u userService) Login(ctx context.Context, auth dto.Auth) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.dbResponseTime)
	defer cancel()

	id, hashedPassword, err := u.userRepo.GetSecure(ctx, auth.Email)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return 0, err
	}

	isCompare := utils.ComparePassword(hashedPassword, auth.Password)
	if !isCompare {
		u.logger.Error().Msg(errs.InvalidPassword.Error())
		return 0, errs.InvalidPassword
	}

	u.logger.Info().Msg(log.Normalizer(log.AuthorizeUser, auth.Email))

	return id, nil
}
