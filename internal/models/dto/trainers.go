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
	ID    int    `json:"id" validate:"required,min=1"`
	Email string `json:"email" validate:"required,email"`
}
