package domain

import (
	"gopkg.in/guregu/null.v3"
)

type TrainerBase struct {
	FirstName  string      `db:"first_name"`
	LastName   string      `db:"last_name"`
	Age        int         `db:"age"`
	Sex        int         `db:"sex"`
	Experience int         `db:"experience"`
	Quote      null.String `db:"quote"`
}

type TrainerCreate struct {
	TrainerBase
	Email    string
	Password string
}

type TrainerUpdate struct {
	TrainerBase
	ID    int
	Email string
}

type TrainerCover struct {
	TrainerBase
	ID              int         `db:"id"`
	PhotoUrl        null.String `db:"photo_url"`
	Roles           []Base      `json:"roles"`
	Specializations []Base      `json:"specializations"`
}

type TrainerCoverPagination struct {
	Trainers []TrainerCover
	Cursor   int
}

type Trainer struct {
	TrainerCover
	Services     []Service
	Achievements []BaseStatus
	Email        string
}

type ServiceBase struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	ProfileAccess bool   `json:"profile_access"`
}

type ServiceCreate struct {
	ServiceBase
	TrainerID int
}

type ServiceUpdate struct {
	ServiceBase
	ID int `json:"id"`
}

type Service struct {
	ServiceUpdate
}
