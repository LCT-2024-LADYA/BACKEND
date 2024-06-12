package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/errs"
	"BACKEND/internal/services"
	"BACKEND/pkg/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ServiceHandler struct {
	service   services.Base
	converter converters.BaseConverter
}

func InitServiceHandler(
	service services.Base,
) *ServiceHandler {
	return &ServiceHandler{
		service:   service,
		converter: converters.InitBaseConverter(),
	}
}

// GetServiceByID
// @Summary Get Service by ID
// @Description Get a service by its ID
// @Tags Services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} dto.BasePrice "Service details"
// @Failure 400 {object} responses.MessageResponse "Invalid service ID"
// @Failure 404 {object} responses.MessageResponse "Service not found"
// @Failure 500 "Internal server error"
// @Router /api/service/{id} [get]
func (s *ServiceHandler) GetServiceByID(c *gin.Context) {
	serviceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	service, err := s.service.GetServiceByID(c.Request.Context(), serviceID)
	if err != nil {
		if errors.Is(err, errs.ErrNoService) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, responses.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}
