package domain

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type UserTrainerServiceCreate struct {
	UserID    int
	TrainerID int
	ServiceID int
}

type ServiceUser struct {
	UserTrainerServiceCreate
	Service        BasePrice
	User           UserCover
	ID             int
	IsPayed        bool
	TrainerConfirm null.Bool
	UserConfirm    null.Bool
}

type ServiceUserPagination struct {
	Services []ServiceUser
	Cursor   int
}

type ServiceTrainer struct {
	UserTrainerServiceCreate
	Service        BasePrice
	Trainer        TrainerCover
	ID             int
	IsPayed        bool
	TrainerConfirm null.Bool
	UserConfirm    null.Bool
}

type ServiceTrainerPagination struct {
	Services []ServiceTrainer
	Cursor   int
}

type ScheduleService struct {
	ScheduleID int
	Date       time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
}

type ScheduleServiceUser struct {
	ServiceUser
	ScheduleService
}
