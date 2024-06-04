package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type TrainerConverter interface {
	TrainerBaseDTOToDomain(trainer dto.TrainerBase) domain.TrainerBase
	TrainerCreateDTOToDomain(trainer dto.TrainerCreate) domain.TrainerCreate
}

type trainerConverter struct{}

func InitTrainerConverter() TrainerConverter {
	return &trainerConverter{}
}

func (t trainerConverter) TrainerBaseDTOToDomain(trainer dto.TrainerBase) domain.TrainerBase {
	return domain.TrainerBase{
		FirstName:  trainer.FirstName,
		LastName:   trainer.LastName,
		Age:        trainer.Age,
		Sex:        trainer.Sex,
		Experience: trainer.Experience,
		Quote:      getNullString(trainer.Quote),
	}
}

func (t trainerConverter) TrainerCreateDTOToDomain(trainer dto.TrainerCreate) domain.TrainerCreate {
	return domain.TrainerCreate{
		TrainerBase: t.TrainerBaseDTOToDomain(trainer.TrainerBase),
		Email:       trainer.Email,
		Password:    trainer.Password,
	}
}
