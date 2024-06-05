package repository

import (
	"BACKEND/internal/models/domain"
	"context"
	"gopkg.in/guregu/null.v3"
)

type Base interface {
	Create(ctx context.Context, base domain.BaseBase) (int, error)
	Get(ctx context.Context) ([]domain.Base, error)
	Delete(ctx context.Context, baseIDs []int) error
}

type Users interface {
	Create(ctx context.Context, user domain.UserCreate) (int, error)
	GetByID(ctx context.Context, userID int) (domain.User, error)
	GetSecure(ctx context.Context, email string) (int, string, error)
	UpdateMain(ctx context.Context, user domain.UserUpdate) error
	UpdatePhotoUrl(ctx context.Context, userID int, newPhotoUrl null.String) error
}

type Trainers interface {
	Create(ctx context.Context, trainer domain.TrainerCreate) (int, error)
	GetByID(ctx context.Context, trainerID int) (domain.Trainer, error)
	GetSecure(ctx context.Context, email string) (int, string, error)
	UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error
	UpdatePhotoUrl(ctx context.Context, trainerID int, newPhotoUrl null.String) error
	UpdateRoles(ctx context.Context, trainerID int, roleIDs []int) error
	UpdateSpecializations(ctx context.Context, trainerID int, specializationIDs []int) error
	CreateService(ctx context.Context, trainerID int, name string, price int) (int, error)
	UpdateService(ctx context.Context, serviceID int, name string, price int) error
	DeleteService(ctx context.Context, trainerID, serviceID int) error
	CreateAchievement(ctx context.Context, trainerID int, achievement string) (int, error)
	UpdateAchievementStatus(ctx context.Context, achievementID int, status bool) error
	DeleteAchievement(ctx context.Context, trainerID, achievementID int) error
}
