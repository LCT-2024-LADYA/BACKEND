package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type TrainerConverter interface {
	TrainerBaseDTOToDomain(trainer dto.TrainerBase) domain.TrainerBase
	TrainerCreateDTOToDomain(trainer dto.TrainerCreate) domain.TrainerCreate
	TrainerUpdateDTOToDomain(trainer dto.TrainerUpdate, trainerID int) domain.TrainerUpdate

	TrainerBaseDomainToDTO(trainer domain.TrainerBase) dto.TrainerBase
	TrainerCoverDomainToDTO(trainer domain.TrainerCover) dto.TrainerCover
	TrainerCoversDomainToDTO(trainers []domain.TrainerCover) []dto.TrainerCover
	TrainerCoverPaginationDomainToDTO(trainer domain.TrainerCoverPagination) dto.TrainerCoverPagination
	TrainerDomainToDTO(trainer domain.Trainer) dto.Trainer
}

type trainerConverter struct {
	baseConverter BaseConverter
}

func InitTrainerConverter() TrainerConverter {
	return &trainerConverter{
		baseConverter: InitBaseConverter(),
	}
}

// DTO -> Domain

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

func (t trainerConverter) TrainerUpdateDTOToDomain(trainer dto.TrainerUpdate, trainerID int) domain.TrainerUpdate {
	return domain.TrainerUpdate{
		TrainerBase: t.TrainerBaseDTOToDomain(trainer.TrainerBase),
		ID:          trainerID,
		Email:       trainer.Email,
	}
}

// Domain -> DTO

func (t trainerConverter) TrainerBaseDomainToDTO(trainer domain.TrainerBase) dto.TrainerBase {
	return dto.TrainerBase{
		FirstName:  trainer.FirstName,
		LastName:   trainer.LastName,
		Age:        trainer.Age,
		Sex:        trainer.Sex,
		Experience: trainer.Experience,
		Quote:      getStringPointer(trainer.Quote),
	}
}

func (t trainerConverter) TrainerCoverDomainToDTO(trainer domain.TrainerCover) dto.TrainerCover {
	return dto.TrainerCover{
		TrainerBase:     t.TrainerBaseDomainToDTO(trainer.TrainerBase),
		Roles:           t.baseConverter.BasesDomainToDTO(trainer.Roles),
		Specializations: t.baseConverter.BasesDomainToDTO(trainer.Specializations),
		ID:              trainer.ID,
		PhotoUrl:        getStringPointer(trainer.PhotoUrl),
	}
}

func (t trainerConverter) TrainerCoversDomainToDTO(trainers []domain.TrainerCover) []dto.TrainerCover {
	result := make([]dto.TrainerCover, len(trainers))

	for i, trainer := range trainers {
		result[i] = t.TrainerCoverDomainToDTO(trainer)
	}

	return result
}

func (t trainerConverter) TrainerCoverPaginationDomainToDTO(trainer domain.TrainerCoverPagination) dto.TrainerCoverPagination {
	return dto.TrainerCoverPagination{
		Trainers: t.TrainerCoversDomainToDTO(trainer.Trainers),
		Cursor:   trainer.Cursor,
	}
}

func (t trainerConverter) TrainerDomainToDTO(trainer domain.Trainer) dto.Trainer {
	return dto.Trainer{
		TrainerCover: t.TrainerCoverDomainToDTO(trainer.TrainerCover),
		Services:     t.baseConverter.BasesPriceDomainToDTO(trainer.Services),
		Achievements: t.baseConverter.BasesStatusDomainToDTO(trainer.Achievements),
		Email:        trainer.Email,
	}
}
