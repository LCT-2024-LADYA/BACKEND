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

type userService struct {
	userRepo       repository.Users
	converter      converters.UserConverter
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
		converter:      converters.InitUserConverter(),
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

func (u userService) GetByID(ctx context.Context, userID int) (dto.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.dbResponseTime)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return dto.User{}, err
	}

	u.logger.Info().Msg(log.Normalizer(log.GetObject, log.User, user.ID))

	return u.converter.UserDomainToDTO(user), nil
}

func (u userService) UpdateMain(ctx context.Context, user domain.UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, u.dbResponseTime)
	defer cancel()

	err := u.userRepo.UpdateMain(ctx, user)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return err
	}

	u.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.User, user.ID))

	return nil
}

// UpdatePhotoUrl `c *gin.Context` передается во избежаниe выноса бизнес логики с хендлерный слой
func (u userService) UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, userID int) error {
	ctx, cancel := context.WithTimeout(c.Request.Context(), u.dbResponseTime)
	defer cancel()

	// Получение информации о пользователе
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return err
	}

	// Проверка на наличие новой фотографии профиля, если ее нет - удаляем нынешнюю
	if newPhoto != nil {
		key := uuid.New().String()

		uniqueFileName := key + filepath.Ext(newPhoto.Filename)
		filePath := fmt.Sprintf("/static/img/users/profile/%s", uniqueFileName)
		if err := c.SaveUploadedFile(newPhoto, ".."+filePath); err != nil {
			u.logger.Error().Msg("Failed to save uploaded file")
			return err
		}

		if user.UserCover.PhotoUrl.Valid {
			// Обновляем фото
			ctxUpdate, cancelUpdate := context.WithTimeout(c.Request.Context(), u.dbResponseTime)
			defer cancelUpdate()

			if err = u.userRepo.UpdatePhotoUrl(ctxUpdate, userID, null.NewString(filePath, true)); err != nil {
				os.Remove(".." + filePath)
				u.logger.Error().Msg(err.Error())
				return err
			}
			os.Remove(".." + user.UserCover.PhotoUrl.String)
		} else {
			// Устанавливаем новое фото
			ctxDelete, cancelDelete := context.WithTimeout(c.Request.Context(), u.dbResponseTime)
			defer cancelDelete()

			if err := u.userRepo.UpdatePhotoUrl(ctxDelete, userID, null.NewString(filePath, true)); err != nil {
				os.Remove(".." + filePath)
				u.logger.Error().Msg(err.Error())
				return err
			}
		}
	} else {
		if user.UserCover.PhotoUrl.Valid {
			// Удаляем старое фото
			ctxDelete, cancelDelete := context.WithTimeout(c.Request.Context(), u.dbResponseTime)
			defer cancelDelete()

			if err := u.userRepo.UpdatePhotoUrl(ctxDelete, userID, null.NewString("", false)); err != nil {
				u.logger.Error().Msg(err.Error())
				return err
			}
			os.Remove(".." + user.UserCover.PhotoUrl.String)
		} else {
			// Хорошего дня и позитивного настроения
			return nil
		}
	}

	u.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.User, userID))

	return nil
}
