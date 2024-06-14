package services

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/pkg/responses"
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
	"mime/multipart"
	"time"
)

type Base interface {
	Create(ctx context.Context, base domain.BaseBase) (int, error)
	GetByName(ctx context.Context) ([]dto.Base, error)
	GetServiceByID(ctx context.Context, id int) (dto.Service, error)
	Delete(ctx context.Context, baseIDs []int) error
}

type Users interface {
	Register(ctx context.Context, user domain.UserCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
	GetByID(ctx context.Context, userID int) (dto.User, error)
	GetCovers(ctx context.Context, search string, cursor int) (dto.UserCoverPagination, error)
	UpdateMain(ctx context.Context, user domain.UserUpdate) error
	UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, userID int) error
}

type Trainers interface {
	Register(ctx context.Context, trainer domain.TrainerCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
	GetByID(ctx context.Context, trainerID int) (dto.Trainer, error)
	GetCovers(ctx context.Context, filters domain.FiltersTrainerCovers) (dto.TrainerCoverPagination, error)
	UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error
	UpdatePhotoUrl(c *gin.Context, newPhoto *multipart.FileHeader, trainerID int) error
	UpdateRoles(ctx context.Context, trainerID int, roleIDs []int) error
	UpdateSpecializations(ctx context.Context, trainerID int, specializationIDs []int) error
	CreateService(ctx context.Context, service domain.ServiceCreate) (int, error)
	UpdateService(ctx context.Context, service domain.ServiceUpdate) error
	DeleteService(ctx context.Context, trainerID, serviceID int) error
	CreateAchievement(ctx context.Context, trainerID int, achievement string) (int, error)
	UpdateAchievementStatus(ctx context.Context, achievementID int, status bool) error
	DeleteAchievement(ctx context.Context, trainerID, achievementID int) error
}

type UserTrainerServices interface {
	Create(ctx context.Context, service domain.UserTrainerServiceCreate) (int, error)
	Schedule(ctx context.Context, schedule domain.ScheduleService) (int, error)
	GetSchedule(ctx context.Context, month, trainerID int) ([]dto.Schedule, error)
	GetSchedulesByIDs(ctx context.Context, scheduleIDs []int) ([]dto.ScheduleServiceUser, error)
	DeleteScheduled(ctx context.Context, scheduleID int) error
	GetUserServices(ctx context.Context, trainerID, cursor int) (dto.ServiceUserPagination, error)
	GetTrainerServices(ctx context.Context, userID, cursor int) (dto.ServiceTrainerPagination, error)
	UpdateStatus(ctx context.Context, field string, serviceID int, status bool) error
	Delete(ctx context.Context, serviceID int) error
}

type Tokens interface {
	Create(ctx context.Context, userID int, userType string) (responses.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (responses.TokenResponse, error)
}

type Trainings interface {
	CreateExercises(ctx context.Context, exercises []domain.ExerciseCreateBase) ([]int, error)
	GetExercises(ctx context.Context, search string, cursor int) (dto.ExercisePagination, error)
	CreateTrainingBases(ctx context.Context, trainings []domain.TrainingCreateBase) ([]int, error)
	CreateTraining(ctx context.Context, training domain.TrainingCreate) (int, error)
	SetExerciseStatus(ctx context.Context, usersTrainingsID, usersExercisesID int, status bool) error
	GetTrainingCovers(ctx context.Context, search string, userID null.Int, cursor int) (dto.TrainingCoverPagination, error)
	GetTraining(ctx context.Context, trainingID int) (dto.Training, error)
	GetScheduleTrainings(ctx context.Context, userTrainingIDs []int) ([]dto.UserTraining, error)
	ScheduleTraining(ctx context.Context, training domain.ScheduleTraining) (int, []int, error)
	GetSchedule(ctx context.Context, month, userID int) ([]dto.Schedule, error)
	DeleteUserTraining(ctx context.Context, trainingID int) error
	DeleteScheduledTraining(ctx context.Context, userTrainingID int) error
}

type Chat interface {
	CreateMessage(ctx context.Context, message domain.MessageCreate) (int, time.Time, error)
	GetUserChats(ctx context.Context, userID int, search string) ([]dto.Chat, error)
	GetTrainerChats(ctx context.Context, trainerID int, search string) ([]dto.Chat, error)
	GetChatMessage(ctx context.Context, userID, trainerID, cursor int) (dto.MessagePagination, error)
}
