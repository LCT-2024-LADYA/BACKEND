package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/delivery/middleware"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/internal/services"
	"BACKEND/pkg/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserTrainerServiceHandler struct {
	service   services.UserTrainerServices
	converter converters.ServicesConverter
}

func InitUserTrainerServiceHandler(
	service services.UserTrainerServices,
) *UserTrainerServiceHandler {
	return &UserTrainerServiceHandler{
		service:   service,
		converter: converters.InitServiceConverter(),
	}
}

// CreateService
// @Summary Create Service
// @Description Create a new service between user and trainer
// @Tags Services
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param service body dto.UserTrainerServiceCreateTrainer true "Service data to create"
// @Success 201 {object} responses.CreatedIDResponse "Service created successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/service [post]
func (s UserTrainerServiceHandler) CreateService(c *gin.Context) {
	var service dto.UserTrainerServiceCreateTrainer

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	id, err := s.service.Create(ctx, s.converter.UserTrainerServiceCreateTrainerDTOToDomain(service, trainerID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// ScheduleService
// @Summary Schedule Service
// @Description Schedule a new service between user and trainer
// @Tags Services
// @Accept json
// @Produce json
// @Param schedule body dto.ScheduleService true "Schedule data to create"
// @Success 201 {object} responses.CreatedIDResponse "Schedule created successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid body provided"
// @Failure 500 "Internal server error"
// @Router /api/service/schedule [post]
func (s UserTrainerServiceHandler) ScheduleService(c *gin.Context) {
	var schedule dto.ScheduleService

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	id, err := s.service.Schedule(ctx, s.converter.ScheduleServiceDTOToDomain(schedule))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// GetSchedule
// @Summary Get Schedule
// @Description Get schedule for the specified month
// @Tags Services
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param month path int true "Month (1-12)"
// @Success 200 {object} []dto.TrainingSchedule "Return schedule for the month"
// @Failure 400 {object} responses.MessageResponse "Invalid month or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/service/schedule/{month} [get]
func (s UserTrainerServiceHandler) GetSchedule(c *gin.Context) {
	monthStr := c.Param("month")
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	schedule, err := s.service.GetSchedule(ctx, month, trainerID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetSchedulesByIDs
// @Summary Get Schedules by IDs
// @Description Get schedules by an array of schedule IDs
// @Tags Services
// @Accept json
// @Produce json
// @Param schedule_ids query []int true "Array of schedule IDs"
// @Success 200 {object} []dto.ScheduleServiceUser "Return schedules for the given IDs"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters provided"
// @Failure 500 "Internal server error"
// @Router /api/service/schedule [get]
func (s UserTrainerServiceHandler) GetSchedulesByIDs(c *gin.Context) {
	scheduleIDsStr := c.QueryArray("schedule_ids")
	if len(scheduleIDsStr) == 0 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	var scheduleIDs []int
	for _, idStr := range scheduleIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
			return
		}
		scheduleIDs = append(scheduleIDs, id)
	}

	ctx := c.Request.Context()

	schedules, err := s.service.GetSchedulesByIDs(ctx, scheduleIDs)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// DeleteScheduled
// @Summary Delete Scheduled Service
// @Description Delete a scheduled service by its ID
// @Tags Services
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 200 "Schedule deleted successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid schedule ID provided"
// @Failure 500 "Internal server error"
// @Router /api/service/schedule/{schedule_id} [delete]
func (s UserTrainerServiceHandler) DeleteScheduled(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	err = s.service.DeleteScheduled(ctx, scheduleID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// GetTrainerServices
// @Summary Get Trainer Services
// @Description Get services for a trainer with pagination
// @Tags Services
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.ServiceUserPagination "List of user services with pagination"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/service/trainer [get]
func (s UserTrainerServiceHandler) GetTrainerServices(c *gin.Context) {
	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	userServices, err := s.service.GetUserServices(ctx, trainerID, cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userServices)
}

// GetUserServices
// @Summary Get User Services
// @Description Get services for a user with pagination
// @Tags Services
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.ServiceTrainerPagination "List of trainer services with pagination"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/service/user [get]
func (s UserTrainerServiceHandler) GetUserServices(c *gin.Context) {
	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	trainerServices, err := s.service.GetTrainerServices(ctx, userID, cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, trainerServices)
}

// UpdateStatus
// @Summary Update Service Status
// @Description Update the status of a service
// @Description type = 1 - payment status
// @Description type = 2 - trainer confirm status
// @Description type = 3 - user confirm status
// @Tags Services
// @Accept json
// @Produce json
// @Param service_id path int true "Service ID"
// @Param status body dto.UpdateStatusService true "Status data to update"
// @Success 200 "Status updated successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/service/status/{service_id} [put]
func (s UserTrainerServiceHandler) UpdateStatus(c *gin.Context) {
	serviceIDStr := c.Param("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	var status dto.UpdateStatusService
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if status.Type < 1 || status.Type > 3 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err = s.service.UpdateStatus(ctx, repository.UsersTrainersServicesField[status.Type], serviceID, status.Status)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteService
// @Summary Delete Service
// @Description Delete a service by its ID
// @Tags Services
// @Accept json
// @Produce json
// @Param service_id path int true "Service ID"
// @Success 200 "Service deleted successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid service ID provided"
// @Failure 500 "Internal server error"
// @Router /api/service/{service_id} [delete]
func (s UserTrainerServiceHandler) DeleteService(c *gin.Context) {
	serviceIDStr := c.Param("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	err = s.service.Delete(ctx, serviceID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
