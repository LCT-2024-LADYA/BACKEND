package services

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"context"
	"github.com/rs/zerolog"
	"time"
)

type usersTrainersServicesService struct {
	serviceRepo    repository.UsersTrainersServices
	converter      converters.ServicesConverter
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitUsersTrainersServicesService(
	serviceRepo repository.UsersTrainersServices,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) UserTrainerServices {
	return &usersTrainersServicesService{
		serviceRepo:    serviceRepo,
		converter:      converters.InitServiceConverter(),
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (s usersTrainersServicesService) Create(ctx context.Context, service domain.UserTrainerServiceCreate) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	createdID, err := s.serviceRepo.Create(ctx, service)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return 0, err
	}

	s.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Service, createdID))

	return createdID, nil
}

func (s usersTrainersServicesService) Schedule(ctx context.Context, schedule domain.ScheduleService) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	createdID, err := s.serviceRepo.Schedule(ctx, schedule)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return 0, err
	}

	s.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Service, createdID))

	return createdID, nil
}

func (s usersTrainersServicesService) GetSchedule(ctx context.Context, month, trainerID int) ([]dto.Schedule, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	schedules, err := s.serviceRepo.GetSchedule(ctx, month, trainerID)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return []dto.Schedule{}, err
	}

	s.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Service))

	return s.converter.SchedulesDomainToDTO(schedules), nil
}

func (s usersTrainersServicesService) GetSchedulesByIDs(ctx context.Context, scheduleIDs []int) ([]dto.ScheduleServiceUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	schedules, err := s.serviceRepo.GetSchedulesByIDs(ctx, scheduleIDs)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return []dto.ScheduleServiceUser{}, err
	}

	s.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Service))

	return s.converter.SchedulesServiceUserDomainToDTO(schedules), nil
}

func (s usersTrainersServicesService) DeleteScheduled(ctx context.Context, scheduleID int) error {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	err := s.serviceRepo.DeleteScheduled(ctx, scheduleID)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	s.logger.Info().Msg(log.Normalizer(log.DeleteObject, log.Service, scheduleID))

	return nil
}

func (s usersTrainersServicesService) GetUserServices(ctx context.Context, trainerID, cursor int) (dto.ServiceUserPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	services, err := s.serviceRepo.GetUserServices(ctx, trainerID, cursor)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return dto.ServiceUserPagination{}, err
	}

	s.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Service))

	return s.converter.ServiceUserPaginationDomainToDTO(services), nil
}

func (s usersTrainersServicesService) GetTrainerServices(ctx context.Context, userID, cursor int) (dto.ServiceTrainerPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	services, err := s.serviceRepo.GetTrainerServices(ctx, userID, cursor)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return dto.ServiceTrainerPagination{}, err
	}

	s.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Trainer))

	return s.converter.ServiceTrainerPaginationDomainToDTO(services), nil
}

func (s usersTrainersServicesService) UpdateStatus(ctx context.Context, field string, serviceID int, status bool) error {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	err := s.serviceRepo.UpdateStatus(ctx, field, serviceID, status)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	s.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Trainer))

	return nil
}

func (s usersTrainersServicesService) Delete(ctx context.Context, serviceID int) error {
	ctx, cancel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cancel()

	err := s.serviceRepo.Delete(ctx, serviceID)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	s.logger.Info().Msg(log.Normalizer(log.DeleteObject, log.Service, serviceID))

	return nil
}
