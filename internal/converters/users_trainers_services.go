package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type ServicesConverter interface {
	UserTrainerServiceCreateTrainerDTOToDomain(service dto.UserTrainerServiceCreateTrainer, trainerID int) domain.UserTrainerServiceCreate
	ScheduleServiceDTOToDomain(schedule dto.ScheduleService) domain.ScheduleService

	UserTrainerServiceCreateDomainToDTO(service domain.UserTrainerServiceCreate) dto.UserTrainerServiceCreate
	ServiceUserDomainToDTO(service domain.ServiceUser) dto.ServiceUser
	ServicesUserDomainToDTO(services []domain.ServiceUser) []dto.ServiceUser
	ServiceUserPaginationDomainToDTO(service domain.ServiceUserPagination) dto.ServiceUserPagination
	ServiceTrainerDomainToDTO(service domain.ServiceTrainer) dto.ServiceTrainer
	ServicesTrainerDomainToDTO(services []domain.ServiceTrainer) []dto.ServiceTrainer
	ServiceTrainerPaginationDomainToDTO(service domain.ServiceTrainerPagination) dto.ServiceTrainerPagination
	ScheduleDomainToDTO(schedule domain.TrainingSchedule) dto.TrainingSchedule
	SchedulesDomainToDTO(schedules []domain.TrainingSchedule) []dto.TrainingSchedule
	ScheduleServiceDomainToDTO(schedule domain.ScheduleService) dto.ScheduleService
	ScheduleServiceUserDomainToDTO(schedule domain.ScheduleServiceUser) dto.ScheduleServiceUser
	SchedulesServiceUserDomainToDTO(schedules []domain.ScheduleServiceUser) []dto.ScheduleServiceUser
}

type servicesConverter struct {
	userConverter    UserConverter
	trainerConverter TrainerConverter
}

func InitServiceConverter() ServicesConverter {
	return &servicesConverter{
		userConverter:    InitUserConverter(),
		trainerConverter: InitTrainerConverter(),
	}
}

func (s servicesConverter) UserTrainerServiceCreateTrainerDTOToDomain(service dto.UserTrainerServiceCreateTrainer, trainerID int) domain.UserTrainerServiceCreate {
	return domain.UserTrainerServiceCreate{
		UserID:    service.UserID,
		TrainerID: trainerID,
		ServiceID: service.ServiceID,
	}
}

func (s servicesConverter) ScheduleServiceDTOToDomain(schedule dto.ScheduleService) domain.ScheduleService {
	return domain.ScheduleService{
		ScheduleID: schedule.ScheduleID,
		Date:       schedule.Date,
		TimeStart:  schedule.TimeStart,
		TimeEnd:    schedule.TimeEnd,
	}
}

func (s servicesConverter) UserTrainerServiceCreateDomainToDTO(service domain.UserTrainerServiceCreate) dto.UserTrainerServiceCreate {
	return dto.UserTrainerServiceCreate{
		UserID:    service.UserID,
		TrainerID: service.TrainerID,
		ServiceID: service.ServiceID,
	}
}

func (s servicesConverter) ServiceUserDomainToDTO(service domain.ServiceUser) dto.ServiceUser {
	return dto.ServiceUser{
		UserTrainerServiceCreate: s.UserTrainerServiceCreateDomainToDTO(service.UserTrainerServiceCreate),
		Service:                  s.trainerConverter.ServiceDomainToDTO(service.Service),
		User:                     s.userConverter.UserCoverDomainToDTO(service.User),
		ID:                       service.ID,
		IsPayed:                  service.IsPayed,
		TrainerConfirm:           getBoolPointer(service.TrainerConfirm),
		UserConfirm:              getBoolPointer(service.UserConfirm),
	}
}

func (s servicesConverter) ServicesUserDomainToDTO(services []domain.ServiceUser) []dto.ServiceUser {
	result := make([]dto.ServiceUser, len(services))

	for i, service := range services {
		result[i] = s.ServiceUserDomainToDTO(service)
	}

	return result
}

func (s servicesConverter) ServiceUserPaginationDomainToDTO(service domain.ServiceUserPagination) dto.ServiceUserPagination {
	return dto.ServiceUserPagination{
		Services: s.ServicesUserDomainToDTO(service.Services),
		Cursor:   service.Cursor,
	}
}

func (s servicesConverter) ServiceTrainerDomainToDTO(service domain.ServiceTrainer) dto.ServiceTrainer {
	return dto.ServiceTrainer{
		UserTrainerServiceCreate: s.UserTrainerServiceCreateDomainToDTO(service.UserTrainerServiceCreate),
		Service:                  s.trainerConverter.ServiceDomainToDTO(service.Service),
		Trainer:                  s.trainerConverter.TrainerCoverDomainToDTO(service.Trainer),
		ID:                       service.ID,
		IsPayed:                  service.IsPayed,
		TrainerConfirm:           getBoolPointer(service.TrainerConfirm),
		UserConfirm:              getBoolPointer(service.UserConfirm),
	}
}

func (s servicesConverter) ServicesTrainerDomainToDTO(services []domain.ServiceTrainer) []dto.ServiceTrainer {
	result := make([]dto.ServiceTrainer, len(services))

	for i, service := range services {
		result[i] = s.ServiceTrainerDomainToDTO(service)
	}

	return result
}

func (s servicesConverter) ServiceTrainerPaginationDomainToDTO(service domain.ServiceTrainerPagination) dto.ServiceTrainerPagination {
	return dto.ServiceTrainerPagination{
		Services: s.ServicesTrainerDomainToDTO(service.Services),
		Cursor:   service.Cursor,
	}
}

func (s servicesConverter) ScheduleDomainToDTO(schedule domain.TrainingSchedule) dto.TrainingSchedule {
	return dto.TrainingSchedule{
		Date:        schedule.Date,
		TrainingIDs: schedule.TrainingIDs,
	}
}

func (s servicesConverter) SchedulesDomainToDTO(schedules []domain.TrainingSchedule) []dto.TrainingSchedule {
	result := make([]dto.TrainingSchedule, len(schedules))

	for i, schedule := range schedules {
		result[i] = s.ScheduleDomainToDTO(schedule)
	}

	return result
}

func (s servicesConverter) ScheduleServiceDomainToDTO(schedule domain.ScheduleService) dto.ScheduleService {
	return dto.ScheduleService{
		ScheduleID: schedule.ScheduleID,
		Date:       schedule.Date,
		TimeStart:  schedule.TimeStart,
		TimeEnd:    schedule.TimeEnd,
	}
}

func (s servicesConverter) ScheduleServiceUserDomainToDTO(schedule domain.ScheduleServiceUser) dto.ScheduleServiceUser {
	return dto.ScheduleServiceUser{
		ServiceUser:     s.ServiceUserDomainToDTO(schedule.ServiceUser),
		ScheduleService: s.ScheduleServiceDomainToDTO(schedule.ScheduleService),
	}
}

func (s servicesConverter) SchedulesServiceUserDomainToDTO(schedules []domain.ScheduleServiceUser) []dto.ScheduleServiceUser {
	result := make([]dto.ScheduleServiceUser, len(schedules))

	for i, schedule := range schedules {
		result[i] = s.ScheduleServiceUserDomainToDTO(schedule)
	}

	return result
}
