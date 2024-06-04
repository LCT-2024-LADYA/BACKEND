package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type UserConverter interface {
	UserBaseDTOToDomain(user dto.UserBase) domain.UserBase
	UserCreateDTOToDomain(user dto.UserCreate) domain.UserCreate
}

type userConverter struct{}

func InitUserConverter() UserConverter {
	return &userConverter{}
}

func (u userConverter) UserBaseDTOToDomain(user dto.UserBase) domain.UserBase {
	return domain.UserBase{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		Sex:       user.Sex,
	}
}

func (u userConverter) UserCreateDTOToDomain(user dto.UserCreate) domain.UserCreate {
	return domain.UserCreate{
		UserBase: u.UserBaseDTOToDomain(user.UserBase),
		Email:    user.Email,
		Password: user.Password,
	}
}
