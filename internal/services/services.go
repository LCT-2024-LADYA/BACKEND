package services

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/pkg/responses"
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Base interface {
	Create(ctx context.Context, base domain.BaseBase) (int, error)
	GetByName(ctx context.Context) ([]dto.Base, error)
	Delete(ctx context.Context, baseIDs []int) error
}

type Users interface {
	Register(ctx context.Context, user domain.UserCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
	GetByID(ctx context.Context, userID int) (dto.User, error)
	UpdateMain(ctx context.Context, user domain.UserUpdate) error
	UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, userID int) error
}

type Trainers interface {
	Register(ctx context.Context, trainer domain.TrainerCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
	GetByID(ctx context.Context, trainerID int) (dto.Trainer, error)
	UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error
	UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, trainerID int) error
	UpdateRoles(ctx context.Context, trainerID int, roleIDs []int) error
	UpdateSpecializations(ctx context.Context, trainerID int, specializationIDs []int) error
	CreateService(ctx context.Context, trainerID int, name string, price int) (int, error)
	UpdateService(ctx context.Context, serviceID int, name string, price int) error
	DeleteService(ctx context.Context, trainerID, serviceID int) error
	CreateAchievement(ctx context.Context, trainerID int, achievement string) (int, error)
	UpdateAchievementStatus(ctx context.Context, achievementID int, status bool) error
	DeleteAchievement(ctx context.Context, trainerID, achievementID int) error
}

type Tokens interface {
	Create(ctx context.Context, userID int, userType string) (responses.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (responses.TokenResponse, error)
}
