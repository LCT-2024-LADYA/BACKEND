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

type TrainerHandler struct {
	service         services.Trainers
	converter       converters.TrainerConverter
	filterConverter converters.FilterConverter
	validate        *validator.Validate
}

func InitTrainerHandler(
	service services.Trainers,
	validate *validator.Validate,
) *TrainerHandler {
	return &TrainerHandler{
		service:         service,
		filterConverter: converters.InitFilterConverter(),
		converter:       converters.InitTrainerConverter(),
		validate:        validate,
	}
}

// Me
// @Summary Get Me
// @Description Get me
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Success 200 {object} dto.Trainer "Return trainer"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal Server Error"
// @Router /api/trainer/me [get]
func (t TrainerHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	trainer, err := t.service.GetByID(ctx, trainerID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUser):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, trainer)
}

// GetProfile
// @Summary Get Profile
// @Description Returned data for profile
// @Tags Trainers
// @Accept json
// @Produce json
// @Param trainer_id path int true "User ID"
// @Success 200 {object} dto.User "Return user"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal Server Error"
// @Router /api/trainer/{trainer_id} [get]
func (t TrainerHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("trainer_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	trainer, err := t.service.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoTrainer):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, trainer)
}

// GetCovers
// @Summary Get Trainer Covers
// @Description Get trainer covers with pagination
// @Tags Trainers
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param cursor query int false "Cursor for pagination"
// @Param role_ids query []int false "Role IDs"
// @Param specialization_ids query []int false "Specialization IDs"
// @Success 200 {object} dto.TrainerCoverPagination "List of trainer covers with pagination"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters"
// @Failure 500 "Internal server error"
// @Router /api/trainer [get]
func (t TrainerHandler) GetCovers(c *gin.Context) {
	var filters dto.FiltersTrainerCovers

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	trainerCovers, err := t.service.GetCovers(c.Request.Context(), t.filterConverter.FilterTrainerDTOToDomain(filters))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, trainerCovers)
}

// UpdateMain
// @Summary Update Trainer's Main Info
// @Description Update trainer's main info by provided data
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param user body dto.TrainerUpdate true "Trainer data to set"
// @Success 200 "Trainer successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/main [put]
func (t TrainerHandler) UpdateMain(c *gin.Context) {
	var trainerUpdate dto.TrainerUpdate

	if err := c.ShouldBindJSON(&trainerUpdate); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := t.validate.Struct(trainerUpdate); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.TrainerUpdate{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	err := t.service.UpdateMain(ctx, t.converter.TrainerUpdateDTOToDomain(trainerUpdate, trainerID))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoTrainer):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePhoto
// @Summary Update Trainer's Photo
// @Description Update trainer's photo
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param photo formData file false "New user photo with type jpeg/jpg/png/svg under 2MB"
// @Success 200 "Trainer successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid photo or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/photo [put]
func (t TrainerHandler) UpdatePhoto(c *gin.Context) {
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

	trainerID := c.GetInt(middleware.UserID)

	err := t.service.UpdatePhotoUrl(c, file, trainerID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoTrainer):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdateRoles
// @Summary Update Trainer's Roles
// @Description Update trainer's roles by provided data
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param roles body []int true "Role IDs to set"
// @Success 200 "Roles successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/roles [put]
func (t TrainerHandler) UpdateRoles(c *gin.Context) {
	var roleIDs []int

	if err := c.ShouldBindJSON(&roleIDs); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	err := t.service.UpdateRoles(ctx, trainerID, roleIDs)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoRole), errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdateSpecializations
// @Summary Update Trainer's Specializations
// @Description Update trainer's specializations by provided data
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param specializations body []int true "Specialization IDs to set"
// @Success 200 "Specializations successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/specializations [put]
func (t TrainerHandler) UpdateSpecializations(c *gin.Context) {
	var specializationIDs []int

	if err := c.ShouldBindJSON(&specializationIDs); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	err := t.service.UpdateSpecializations(ctx, trainerID, specializationIDs)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoSpecialization), errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// CreateService
// @Summary Create Trainer's Service
// @Description Create a new service for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param service body dto.ServiceCreate true "Service data to create"
// @Success 201 {object} responses.CreatedIDResponse "Service successfully created"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/service [post]
func (t TrainerHandler) CreateService(c *gin.Context) {
	var service dto.ServiceCreate

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	serviceID, err := t.service.CreateService(ctx, trainerID, service.Name, service.Price)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: serviceID})
}

// UpdateService
// @Summary Update Trainer's Service
// @Description Update an existing service for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param service body dto.ServiceUpdate true "Service data to update"
// @Success 200 "Service successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/service [put]
func (t TrainerHandler) UpdateService(c *gin.Context) {
	var service dto.ServiceUpdate

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err := t.service.UpdateService(ctx, service.ID, service.Name, service.Price)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoService):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteService
// @Summary Delete Trainer's Service
// @Description Delete an existing service for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param service_id path int true "Service ID"
// @Success 200 "Service successfully deleted"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/service/{service_id} [delete]
func (t TrainerHandler) DeleteService(c *gin.Context) {
	serviceIDStr := c.Param("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()
	trainerID := c.GetInt(middleware.UserID)

	err = t.service.DeleteService(ctx, trainerID, serviceID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoService):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// CreateAchievement
// @Summary Create Trainer's Achievement
// @Description Create a new achievement for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param achievement body dto.AchievementCreate true "Achievement data to create"
// @Success 201 {object} responses.CreatedIDResponse "Achievement successfully created"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/achievement [post]
func (t TrainerHandler) CreateAchievement(c *gin.Context) {
	var achievement dto.AchievementCreate

	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	achievementID, err := t.service.CreateAchievement(ctx, trainerID, achievement.Name)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: achievementID})
}

// UpdateAchievementStatus
// @Summary Update Trainer's Achievement Status
// @Description Update the status of an existing achievement for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param achievement_id path int true "Achievement ID"
// @Param status body dto.AchievementStatusUpdate true "Achievement status to update"
// @Success 200 "Achievement status successfully updated"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/achievement/{achievement_id}/status [put]
func (t TrainerHandler) UpdateAchievementStatus(c *gin.Context) {
	achievementIDStr := c.Param("achievement_id")
	achievementID, err := strconv.Atoi(achievementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	var statusUpdate dto.AchievementStatusUpdate
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err = t.service.UpdateAchievementStatus(ctx, achievementID, statusUpdate.Status)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoAchievement):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteAchievement
// @Summary Delete Trainer's Achievement
// @Description Delete an existing achievement for the trainer
// @Tags Trainers
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param achievement_id path int true "Achievement ID"
// @Success 200 "Achievement successfully deleted"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/trainer/achievement/{achievement_id} [delete]
func (t TrainerHandler) DeleteAchievement(c *gin.Context) {
	achievementIDStr := c.Param("achievement_id")
	achievementID, err := strconv.Atoi(achievementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	err = t.service.DeleteAchievement(ctx, trainerID, achievementID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoAchievement):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
