package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type FilterConverter interface {
	FilterTrainerDTOToDomain(filter dto.FiltersTrainerCovers) domain.FiltersTrainerCovers
}

type filterConverter struct {
}

func InitFilterConverter() FilterConverter {
	return &filterConverter{}
}

func (f filterConverter) FilterTrainerDTOToDomain(filter dto.FiltersTrainerCovers) domain.FiltersTrainerCovers {
	return domain.FiltersTrainerCovers{
		Search:            filter.Search,
		Cursor:            filter.Cursor,
		RoleIDs:           filter.RoleIDs,
		SpecializationIDs: filter.SpecializationIDs,
	}
}
