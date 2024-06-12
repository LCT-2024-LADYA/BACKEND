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

type baseService struct {
	baseRepo       repository.Base
	converter      converters.BaseConverter
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitBaseService(
	baseRepo repository.Base,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Base {
	return &baseService{
		baseRepo:       baseRepo,
		converter:      converters.InitBaseConverter(),
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (b baseService) Create(ctx context.Context, base domain.BaseBase) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, b.dbResponseTime)
	defer cancel()

	createdID, err := b.baseRepo.Create(ctx, base)
	if err != nil {
		b.logger.Error().Msg(err.Error())
		return 0, err
	}

	b.logger.Info().Msg(log.Normalizer(log.CreateObject, b.baseRepo.GetTable(), createdID))

	return createdID, nil
}

func (b baseService) GetByName(ctx context.Context) ([]dto.Base, error) {
	ctx, cancel := context.WithTimeout(ctx, b.dbResponseTime)
	defer cancel()

	bases, err := b.baseRepo.Get(ctx)
	if err != nil {
		b.logger.Error().Msg(err.Error())
		return []dto.Base{}, err
	}

	b.logger.Info().Msg(log.Normalizer(log.GetObjects, b.baseRepo.GetTable()))

	return b.converter.BasesDomainToDTO(bases), nil
}

func (b baseService) GetServiceByID(ctx context.Context, id int) (dto.BasePrice, error) {
	ctx, cancel := context.WithTimeout(ctx, b.dbResponseTime)
	defer cancel()

	service, err := b.baseRepo.GetServiceByID(ctx, id)
	if err != nil {
		b.logger.Error().Msg(err.Error())
		return dto.BasePrice{}, err
	}

	b.logger.Info().Msg(log.Normalizer(log.GetObject, "service", id))

	return b.converter.BasePriceDomainToDTO(service), nil
}

func (b baseService) Delete(ctx context.Context, baseIDs []int) error {
	ctx, cancel := context.WithTimeout(ctx, b.dbResponseTime)
	defer cancel()

	err := b.baseRepo.Delete(ctx, baseIDs)
	if err != nil {
		b.logger.Error().Msg(err.Error())
		return err
	}

	b.logger.Info().Msg(log.Normalizer(log.DeleteObjects, b.baseRepo.GetTable(), baseIDs))

	return nil
}
