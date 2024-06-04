package handlers

import (
	"BACKEND/internal/services"
	"github.com/go-playground/validator/v10"
)

type TrainerHandler struct {
	service  services.Trainers
	validate *validator.Validate
}

func InitTrainerHandler(
	service services.Trainers,
	validate *validator.Validate,
) *TrainerHandler {
	return &TrainerHandler{
		service:  service,
		validate: validate,
	}
}
