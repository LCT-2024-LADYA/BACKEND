package services

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

type userService struct {
	userRepo       repository.Users
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitUserService(
	userRepo repository.Users,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Users {
	return &userService{
		userRepo:       userRepo,
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (u userService) CreateUserIfNotExistByVK(ctx context.Context, user dto.AuthRequest) (int, error) {
	checkExistCtx, cancelExist := context.WithTimeout(ctx, u.dbResponseTime)
	defer cancelExist()

	id, err := u.userRepo.CheckIfExistByVKID(checkExistCtx, user.VKID)
	if err != nil {
		u.logger.Error().Msg(err.Error())
		return 0, err
	}

	if id == -1 {
		userData, err := u.getUserInfoFromVK(user.Token)
		if err != nil {
			u.logger.Error().Msg(err.Error())
			return 0, err
		}
		userData.VKID = user.VKID
		userData.Email = user.Email

		createCtx, cancelCreate := context.WithTimeout(ctx, u.dbResponseTime)
		defer cancelCreate()

		createdID, err := u.userRepo.CreateVK(createCtx, userData)
		if err != nil {
			u.logger.Error().Msg(err.Error())
			return 0, err
		}

		u.logger.Info().Msg(log.Normalizer(log.AuthorizeVK, user.VKID, createdID))

		return createdID, nil
	} else {

		u.logger.Info().Msg(log.Normalizer(log.AuthorizeVK, user.VKID, id))

		return id, nil
	}
}

func (u userService) getUserInfoFromVK(accessToken string) (domain.UserCreateVK, error) {
	fields := "first_name,last_name,bdate,sex"
	apiVersion := "5.131"
	requestURL := fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&fields=%s&v=%s", accessToken, fields, apiVersion)

	resp, err := http.Get(requestURL)
	if err != nil {
		return domain.UserCreateVK{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.UserCreateVK{}, err
	}

	type responseModel struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthDate string `json:"bdate"`
		Sex       int    `json:"sex"`
	}

	var response struct {
		Response []responseModel `json:"response"`
		Error    struct {
			ErrorCode int    `json:"error_code"`
			ErrorMsg  string `json:"error_msg"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return domain.UserCreateVK{}, err
	}

	if response.Error.ErrorCode != 0 {
		return domain.UserCreateVK{}, fmt.Errorf("VK API error %d: %s", response.Error.ErrorCode, response.Error.ErrorMsg)
	}

	if len(response.Response) > 0 {
		var userInfo domain.UserCreateVK

		birthDate, err := time.Parse("2.1.2006", response.Response[0].BirthDate)
		if err == nil {
			now := time.Now()
			age := now.Year() - birthDate.Year()
			if now.YearDay() < birthDate.YearDay() {
				age-- // Уменьшаем на 1, если день рождения в этом году ещё не наступил
			}
			userInfo.Age = age
		} else {
			return domain.UserCreateVK{}, fmt.Errorf("error parsing birthdate: %v", err)
		}

		userInfo.FirstName = response.Response[0].FirstName
		userInfo.LastName = response.Response[0].LastName
		userInfo.Sex = response.Response[0].Sex

		return userInfo, nil
	}

	return domain.UserCreateVK{}, fmt.Errorf("no user data returned")
}
