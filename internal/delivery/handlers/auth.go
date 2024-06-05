package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/config"
	"BACKEND/pkg/responses"
	"BACKEND/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"net/http"
)

type AuthHandler struct {
	userService      services.Users
	trainerService   services.Trainers
	tokenService     services.Tokens
	userConverter    converters.UserConverter
	trainerConverter converters.TrainerConverter
	validate         *validator.Validate
	APIKEY           string
}

func InitAuthHandler(
	userService services.Users,
	trainerService services.Trainers,
	tokenService services.Tokens,
	validate *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		userService:      userService,
		trainerService:   trainerService,
		tokenService:     tokenService,
		userConverter:    converters.InitUserConverter(),
		trainerConverter: converters.InitTrainerConverter(),
		validate:         validate,
		APIKEY:           viper.GetString(config.APIKEY),
	}
}

// RegisterUser
// @Summary User Register
// @Description Register user
// @Tags Authorization
// @Accept json
// @Produce json
// @Param auth body dto.UserCreate true "Register request body"
// @Success 201 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/register/user [post]
func (a AuthHandler) RegisterUser(c *gin.Context) {
	var user dto.UserCreate
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := a.validate.Struct(user); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.Auth{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	id, err := a.userService.Register(ctx, a.userConverter.UserCreateDTOToDomain(user))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
	}

	tokens, err := a.tokenService.Create(ctx, id, utils.User)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

// AuthorizeUser
// @Summary User Authorization
// @Description Authorize user
// @Tags Authorization
// @Accept json
// @Produce json
// @Param auth body dto.Auth true "Authorization request body"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/login/user [post]
func (a AuthHandler) AuthorizeUser(c *gin.Context) {
	var user dto.Auth
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := a.validate.Struct(user); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.Auth{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	id, err := a.userService.Login(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, errs.InvalidEmail), errors.Is(err, errs.InvalidPassword):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	tokens, err := a.tokenService.Create(ctx, id, utils.User)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RegisterTrainer
// @Summary Trainer Register
// @Description Register trainer
// @Tags Authorization
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param auth body dto.TrainerCreate true "Register request body"
// @Success 201 {object} responses.CreatedIDResponse "Return created trainer's id"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/register/trainer [post]
func (a AuthHandler) RegisterTrainer(c *gin.Context) {
	var trainer dto.TrainerCreate
	if err := c.ShouldBindJSON(&trainer); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := a.validate.Struct(trainer); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.Auth{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	id, err := a.trainerService.Register(ctx, a.trainerConverter.TrainerCreateDTOToDomain(trainer))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// AuthorizeTrainer
// @Summary Trainer Authorization
// @Description Authorize trainer
// @Tags Authorization
// @Accept json
// @Produce json
// @Param auth body dto.Auth true "Authorization request body"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/login/trainer [post]
func (a AuthHandler) AuthorizeTrainer(c *gin.Context) {
	var trainer dto.Auth
	if err := c.ShouldBindJSON(&trainer); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := a.validate.Struct(trainer); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.Auth{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	id, err := a.trainerService.Login(ctx, trainer)
	if err != nil {
		switch {
		case errors.Is(err, errs.InvalidEmail), errors.Is(err, errs.InvalidPassword):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	tokens, err := a.tokenService.Create(ctx, id, utils.Trainer)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// AuthorizeAdmin
// @Summary Admin Authorization
// @Description Authorize admin with X-API-KEY
// @Tags Authorization
// @Accept json
// @Produce json
// @Param X-API-KEY header string true "Admin apikey"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 401 {object} responses.MessageResponse "Invalid X-API-KEY"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/login/admin [post]
func (a AuthHandler) AuthorizeAdmin(c *gin.Context) {
	apiKey := c.GetHeader("X-API-KEY")

	if apiKey != a.APIKEY {
		c.JSON(http.StatusUnauthorized, responses.MessageResponse{Message: "Api key is invalid"})
		return
	}

	tokens, err := a.tokenService.Create(c.Request.Context(), 0, utils.Admin)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh
// @Summary Refresh Tokens
// @Description Refreshes the access and refresh tokens using the provided refresh token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param refresh_token query string true "Refresh token"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad query provided"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/refresh [get]
func (a AuthHandler) Refresh(c *gin.Context) {
	refreshToken, ok := c.GetQuery("refresh_token")
	if !ok {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	tokens, err := a.tokenService.Refresh(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, errs.NeedToAuth):
			c.JSON(http.StatusUnauthorized, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}
