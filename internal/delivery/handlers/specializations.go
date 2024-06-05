package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type SpecializationHandler struct {
	service   services.Base
	converter converters.BaseConverter
	validate  *validator.Validate
}

func InitSpecializationHandler(
	service services.Base,
	validate *validator.Validate,
) *SpecializationHandler {
	return &SpecializationHandler{
		service:   service,
		converter: converters.InitBaseConverter(),
		validate:  validate,
	}
}

// CreateSpecialization
// @Summary Create Specialization
// @Description Create a new specialization
// @Tags Specializations
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param specialization body dto.BaseBase true "Specializations data to create"
// @Success 201 {object} responses.CreatedIDResponse "Specializations successfully created"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/specialization [post]
func (s SpecializationHandler) CreateSpecialization(c *gin.Context) {
	var specializations dto.BaseBase

	if err := c.ShouldBindJSON(&specializations); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := s.validate.Struct(specializations); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.BaseBase{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	specializationsID, err := s.service.Create(ctx, s.converter.BaseBaseDTOToDomain(specializations))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: specializationsID})
}

// GetSpecializations
// @Summary Get Specializations
// @Description Get specializations
// @Tags Specializations
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Base "Return specializations"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/specialization [get]
func (s SpecializationHandler) GetSpecializations(c *gin.Context) {
	ctx := c.Request.Context()

	specializations, err := s.service.GetByName(ctx)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, specializations)
}

// DeleteSpecializations
// @Summary Delete Specializations
// @Description Delete specializations by provided IDs
// @Tags Specializations
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param specialization_ids body []int true "Specialization IDs to delete"
// @Success 200 "Specializations successfully deleted"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/specialization [delete]
func (s SpecializationHandler) DeleteSpecializations(c *gin.Context) {
	var specializationIDs []int

	if err := c.ShouldBindJSON(&specializationIDs); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err := s.service.Delete(ctx, specializationIDs)
	if err != nil {
		if err.Error() == "Один из переданных объектов не существует" {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
