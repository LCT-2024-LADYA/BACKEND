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

	GetTable() string
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

type Trainings interface {
	CreateExercises(ctx context.Context, exercises []domain.ExerciseCreateBase) ([]int, error)
	GetExercises(ctx context.Context, search string, cursor int) (domain.ExercisePagination, error)
	CreateTrainingBases(ctx context.Context, trainings []domain.TrainingCreateBase) ([]int, error)
	CreateTraining(ctx context.Context, training domain.TrainingCreate) (int, []int, error)
	SetExerciseStatus(ctx context.Context, usersTrainingsID, usersExercisesID int, status bool) error
	GetTrainingCovers(ctx context.Context, search string, userID null.Int, cursor int) (domain.TrainingCoverPagination, error)
	GetTraining(ctx context.Context, trainingID int) (domain.Training, error)
	GetTrainingsDate(ctx context.Context, userTrainingIDs []int) ([]domain.TrainingDate, error)
	ScheduleTraining(ctx context.Context, training domain.ScheduleTraining) (int, []int, error)
	GetSchedule(ctx context.Context, month, userID int) ([]domain.Schedule, error)
}
