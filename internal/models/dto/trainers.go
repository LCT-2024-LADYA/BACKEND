package dto

type TrainerBase struct {
	FirstName  string  `json:"first_name" validate:"required,min=2,max=50"`
	LastName   string  `json:"last_name" validate:"required,min=2,max=50"`
	Age        int     `json:"age" validate:"required,min=18,max=150"`
	Sex        int     `json:"sex" validate:"required,oneof=1 2"`
	Experience int     `json:"experience" validate:"required,min=0,max=50"`
	Quote      *string `json:"quote" validate:"omitempty,max=100"`
}

type TrainerCreate struct {
	TrainerBase
	Auth
}

type TrainerUpdate struct {
	TrainerBase
	Email string `json:"email" validate:"required,email"`
}

type TrainerCover struct {
	TrainerBase
	Roles           []Base  `json:"roles"`
	Specializations []Base  `json:"specializations"`
	ID              int     `json:"id"`
	PhotoUrl        *string `json:"photo_url"`
}

type TrainerCoverPagination struct {
	Trainers []TrainerCover `json:"objects"`
	Cursor   int            `json:"cursor"`
}

type Trainer struct {
	TrainerCover
	Services     []Service    `json:"services"`
	Achievements []BaseStatus `json:"achievements"`
	Email        string       `json:"email"`
}

type ServiceBase struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	ProfileAccess bool   `json:"profile_access"`
}

type ServiceCreate struct {
	ServiceBase
}

type ServiceUpdate struct {
	ServiceBase
	ID int `json:"id"`
}

type Service struct {
	ServiceUpdate
}

type AchievementCreate struct {
	Name string `json:"name" validate:"required,min=2,max=150"`
}

type AchievementStatusUpdate struct {
	Status bool `json:"status"`
}
