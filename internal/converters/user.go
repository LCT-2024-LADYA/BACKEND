package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type UserConverter interface {
	UserBaseDTOToDomain(user dto.UserBase) domain.UserBase
	UserCreateDTOToDomain(user dto.UserCreate) domain.UserCreate
	UserUpdateDTOToDomain(user dto.UserUpdate, userID int) domain.UserUpdate

	UserBaseDomainToDTO(user domain.UserBase) dto.UserBase
	UserCoverDomainToDTO(user domain.UserCover) dto.UserCover
	UserDomainToDTO(user domain.User) dto.User
}

type userConverter struct{}

func InitUserConverter() UserConverter {
	return &userConverter{}
}

// DTO -> Domain

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

func (u userConverter) UserUpdateDTOToDomain(user dto.UserUpdate, userID int) domain.UserUpdate {
	return domain.UserUpdate{
		UserBase: u.UserBaseDTOToDomain(user.UserBase),
		ID:       userID,
		Email:    user.Email,
	}
}

// Domain -> DTO

func (u userConverter) UserBaseDomainToDTO(user domain.UserBase) dto.UserBase {
	return dto.UserBase{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		Sex:       user.Sex,
	}
}

func (u userConverter) UserCoverDomainToDTO(user domain.UserCover) dto.UserCover {
	return dto.UserCover{
		UserBase: u.UserBaseDomainToDTO(user.UserBase),
		ID:       user.ID,
		PhotoUrl: getStringPointer(user.PhotoUrl),
	}
}

func (u userConverter) UserDomainToDTO(user domain.User) dto.User {
	return dto.User{
		UserCover: u.UserCoverDomainToDTO(user.UserCover),
		Email:     user.Email,
	}
}
