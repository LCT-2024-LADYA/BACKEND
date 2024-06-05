package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/delivery/middleware"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service   services.Users
	converter converters.UserConverter
	validate  *validator.Validate
}

func InitUserHandler(
	service services.Users,
	validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		service:   service,
		converter: converters.InitUserConverter(),
		validate:  validate,
	}
}

// Me
// @Summary Get Me
// @Description Get me
// @Tags Users
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Success 200 {object} dto.User "Return user"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal Server Error"
// @Router /api/user/me [get]
func (u UserHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	user, err := u.service.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUser):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetProfile
// @Summary Get Profile
// @Description Returned data for profile
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.User "Return user"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 500 "Internal Server Error"
// @Router /api/user/{user_id} [get]
func (u UserHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	user, err := u.service.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUser):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMain
// @Summary Update User's Main Info
// @Description Update user's main info by provided data
// @Tags Users
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param user body dto.UserUpdate true "User data to set"
// @Success 200 "User successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/user/main [put]
func (u UserHandler) UpdateMain(c *gin.Context) {
	var newUserMain dto.UserUpdate

	if err := c.ShouldBindJSON(&newUserMain); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := u.validate.Struct(newUserMain); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.UserUpdate{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	err := u.service.UpdateMain(ctx, u.converter.UserUpdateDTOToDomain(newUserMain, userID))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUser):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePhoto
// @Summary Update User's Photo
// @Description Update user's photo
// @Tags Users
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param photo formData file false "New user photo with type jpeg/jpg/png/svg under 2MB"
// @Success 200 "User successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid photo or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/user/photo [put]
func (u UserHandler) UpdatePhoto(c *gin.Context) {
	file, _ := c.FormFile("photo")

	if file != nil {
		// Проверка на ограничение по размеру файла на 2МБ
		err := c.Request.ParseMultipartForm(2 << 20)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: "Photo size is bigger than 2MB"})
			return
		}

		// Проверка на допустимый тип `Content-Type` и расширение
		if !validators.ValidateFileTypeExtension(file) {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: "Photo bad type or extension"})
			return
		}
	}

	userID := c.GetInt(middleware.UserID)

	err := u.service.UpdatePhotoUrl(c, file, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUser):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
