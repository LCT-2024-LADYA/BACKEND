package services

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"BACKEND/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/guregu/null.v3"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type trainerService struct {
	trainerRepo    repository.Trainers
	converter      converters.TrainerConverter
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
		converter:      converters.InitTrainerConverter(),
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

func (t trainerService) GetByID(ctx context.Context, trainerID int) (dto.Trainer, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	trainer, err := t.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.Trainer{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObject, log.Trainer, trainer.ID))

	return t.converter.TrainerDomainToDTO(trainer), nil
}

func (t trainerService) GetCovers(ctx context.Context, filters domain.FiltersTrainerCovers) (dto.TrainerCoverPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	trainers, err := t.trainerRepo.GetCovers(ctx, filters)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.TrainerCoverPagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Trainer))

	return t.converter.TrainerCoverPaginationDomainToDTO(trainers), nil
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

// UpdatePhotoUrl `c *gin.Context` передается во избежаниe выноса бизнес логики с хендлерный слой
func (t trainerService) UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, trainerID int) error {
	ctx, cancel := context.WithTimeout(c.Request.Context(), t.dbResponseTime)
	defer cancel()

	// Получение информации о пользователе
	user, err := t.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	// Проверка на наличие новой фотографии профиля, если ее нет - удаляем нынешнюю
	if newPhoto != nil {
		key := uuid.New().String()

		uniqueFileName := key + filepath.Ext(newPhoto.Filename)
		filePath := fmt.Sprintf("/static/img/trainers/profile/%s", uniqueFileName)
		if err := c.SaveUploadedFile(newPhoto, ".."+filePath); err != nil {
			t.logger.Error().Msg("Failed to save uploaded file")
			return err
		}

		if user.TrainerCover.PhotoUrl.Valid {
			// Обновляем фото
			ctxUpdate, cancelUpdate := context.WithTimeout(c.Request.Context(), t.dbResponseTime)
			defer cancelUpdate()

			if err = t.trainerRepo.UpdatePhotoUrl(ctxUpdate, trainerID, null.NewString(filePath, true)); err != nil {
				os.Remove(".." + filePath)
				t.logger.Error().Msg(err.Error())
				return err
			}
			os.Remove(".." + user.TrainerCover.PhotoUrl.String)
		} else {
			// Устанавливаем новое фото
			ctxDelete, cancelDelete := context.WithTimeout(c.Request.Context(), t.dbResponseTime)
			defer cancelDelete()

			if err := t.trainerRepo.UpdatePhotoUrl(ctxDelete, trainerID, null.NewString(filePath, true)); err != nil {
				os.Remove(".." + filePath)
				t.logger.Error().Msg(err.Error())
				return err
			}
		}
	} else {
		if user.TrainerCover.PhotoUrl.Valid {
			// Удаляем старое фото
			ctxDelete, cancelDelete := context.WithTimeout(c.Request.Context(), t.dbResponseTime)
			defer cancelDelete()

			if err := t.trainerRepo.UpdatePhotoUrl(ctxDelete, trainerID, null.NewString("", false)); err != nil {
				t.logger.Error().Msg(err.Error())
				return err
			}
			os.Remove(".." + user.TrainerCover.PhotoUrl.String)
		} else {
			// Хорошего дня и позитивного настроения
			return nil
		}
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.User, trainerID))

	return nil
}

func (t trainerService) UpdateRoles(ctx context.Context, trainerID int, roleIDs []int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.UpdateRoles(ctx, trainerID, roleIDs)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainerID))

	return nil
}

func (t trainerService) UpdateSpecializations(ctx context.Context, trainerID int, specializationIDs []int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.UpdateSpecializations(ctx, trainerID, specializationIDs)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainerID))

	return nil
}

func (t trainerService) CreateService(ctx context.Context, trainerID int, name string, price int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	createdID, err := t.trainerRepo.CreateService(ctx, trainerID, name, price)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainerID))

	return createdID, nil
}

func (t trainerService) UpdateService(ctx context.Context, serviceID int, name string, price int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.UpdateService(ctx, serviceID, name, price)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, 1505))

	return nil
}

func (t trainerService) DeleteService(ctx context.Context, trainerID, serviceID int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.DeleteService(ctx, trainerID, serviceID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, 1505))

	return nil
}

func (t trainerService) CreateAchievement(ctx context.Context, trainerID int, achievement string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	createdID, err := t.trainerRepo.CreateAchievement(ctx, trainerID, achievement)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainerID))

	return createdID, nil
}

func (t trainerService) UpdateAchievementStatus(ctx context.Context, achievementID int, status bool) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.UpdateAchievementStatus(ctx, achievementID, status)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, 1505))

	return nil
}

func (t trainerService) DeleteAchievement(ctx context.Context, trainerID, achievementID int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainerRepo.DeleteAchievement(ctx, trainerID, achievementID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Trainer, trainerID))

	return nil
}
