package dto

import "time"

type UserTrainerServiceCreateTrainer struct {
	UserID    int `json:"user_id"`
	ServiceID int `json:"service_id"`
}

type UserTrainerServiceCreate struct {
	UserID    int `json:"user_id"`
	TrainerID int `json:"trainer_id"`
	ServiceID int `json:"service_id"`
}

type ServiceUser struct {
	UserTrainerServiceCreate
	Service        BasePrice `json:"service"`
	User           UserCover `json:"user"`
	ID             int       `json:"id"`
	IsPayed        bool      `json:"is_payed"`
	TrainerConfirm *bool     `json:"trainer_confirm"`
	UserConfirm    *bool     `json:"user_confirm"`
}

type ServiceUserPagination struct {
	Services []ServiceUser `json:"objects"`
	Cursor   int           `json:"cursor"`
}

type ServiceTrainer struct {
	UserTrainerServiceCreate
	Service        BasePrice    `json:"service"`
	Trainer        TrainerCover `json:"trainer"`
	ID             int          `json:"id"`
	IsPayed        bool         `json:"is_payed"`
	TrainerConfirm *bool        `json:"trainer_confirm"`
	UserConfirm    *bool        `json:"user_confirm"`
}

type ServiceTrainerPagination struct {
	Services []ServiceTrainer `json:"objects"`
	Cursor   int              `json:"cursor"`
}

type UpdateStatusService struct {
	Status bool `json:"status"`
	Type   int  `json:"type"`
}

type ScheduleService struct {
	ScheduleID int       `json:"schedule_id"`
	Date       time.Time `json:"date"`
	TimeStart  time.Time `json:"time_start"`
	TimeEnd    time.Time `json:"time_end"`
}

type ScheduleServiceUser struct {
	ServiceUser
	ScheduleService
}
