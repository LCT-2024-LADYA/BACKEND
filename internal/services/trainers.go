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

type trainerService struct {
	trainerRepo    repository.Trainers
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitTrainerService(
	trainerRepo repository.Trainers,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Trainers {
	return &trainerService{
		trainerRepo:    trainerRepo,
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (t trainerService) Register(ctx context.Context, trainer domain.TrainerCreate) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	createdID, err := t.trainerRepo.Create(ctx, trainer)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Trainer, createdID))

	return createdID, nil
}

func (t trainerService) Login(ctx context.Context, auth dto.Auth) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	id, hashedPassword, err := t.trainerRepo.GetSecure(ctx, auth.Email)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	isCompare := utils.ComparePassword(hashedPassword, auth.Password)
	if !isCompare {
		t.logger.Error().Msg(errs.InvalidPassword.Error())
		return 0, errs.InvalidPassword
	}

	t.logger.Info().Msg(log.Normalizer(log.AuthorizeTrainer, auth.Email))

	return id, nil
}

func (t trainerService) UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.UpdateMain(ctx, trainer)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainer.ID))

	return nil
}
