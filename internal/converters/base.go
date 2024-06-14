package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type BaseConverter interface {
	BaseBaseDTOToDomain(baseBase dto.BaseBase) domain.BaseBase

	BaseBaseDomainToDTO(baseBase domain.BaseBase) dto.BaseBase
	BaseDomainToDTO(base domain.Base) dto.Base
	BasesDomainToDTO(bases []domain.Base) []dto.Base
	BaseStatusDomainToDTO(base domain.BaseStatus) dto.BaseStatus
	BasesStatusDomainToDTO(bases []domain.BaseStatus) []dto.BaseStatus
	ServiceDomainToDTO(service domain.Service) dto.Service
}

type baseConverter struct{}

func InitBaseConverter() BaseConverter {
	return &baseConverter{}
}

// DTO -> Domain

func (b baseConverter) BaseBaseDTOToDomain(baseBase dto.BaseBase) domain.BaseBase {
	return domain.BaseBase{Name: baseBase.Name}
}

// Domain -> DTO

func (b baseConverter) BaseBaseDomainToDTO(baseBase domain.BaseBase) dto.BaseBase {
	return dto.BaseBase{Name: baseBase.Name}
}

func (b baseConverter) BaseDomainToDTO(base domain.Base) dto.Base {
	return dto.Base{
		ID:       base.ID,
		BaseBase: b.BaseBaseDomainToDTO(base.BaseBase),
	}
}

func (b baseConverter) BasesDomainToDTO(bases []domain.Base) []dto.Base {
	dtoBases := make([]dto.Base, len(bases))

	for i, base := range bases {
		dtoBases[i] = b.BaseDomainToDTO(base)
	}

	return dtoBases
}

func (b baseConverter) BaseStatusDomainToDTO(base domain.BaseStatus) dto.BaseStatus {
	return dto.BaseStatus{
		Base:        b.BaseDomainToDTO(base.Base),
		IsConfirmed: base.IsConfirmed,
	}
}

func (b baseConverter) BasesStatusDomainToDTO(bases []domain.BaseStatus) []dto.BaseStatus {
	dtoBases := make([]dto.BaseStatus, len(bases))

	for i, base := range bases {
		dtoBases[i] = b.BaseStatusDomainToDTO(base)
	}

	return dtoBases
}

func (b baseConverter) ServiceBaseDomainToDTO(service domain.ServiceBase) dto.ServiceBase {
	return dto.ServiceBase{
		Name:          service.Name,
		Price:         service.Price,
		ProfileAccess: service.ProfileAccess,
	}
}

func (b baseConverter) ServiceUpdateDomainToDTO(service domain.ServiceUpdate) dto.ServiceUpdate {
	return dto.ServiceUpdate{
		ServiceBase: b.ServiceBaseDomainToDTO(service.ServiceBase),
		ID:          service.ID,
	}
}

func (b baseConverter) ServiceDomainToDTO(service domain.Service) dto.Service {
	return dto.Service{
		ServiceUpdate: b.ServiceUpdateDomainToDTO(service.ServiceUpdate),
	}
}
