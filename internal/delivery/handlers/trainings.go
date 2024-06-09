package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/delivery/middleware"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/pkg/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
)

type TrainingHandler struct {
	service   services.Trainings
	converter converters.TrainingConverter
}

func InitTrainingsHandler(
	service services.Trainings,
) *TrainingHandler {
	return &TrainingHandler{
		service:   service,
		converter: converters.InitTrainingConverter(),
	}
}

// CreateExercises
// @Summary Create Exercises
// @Description Create multiple exercises
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param exercises body []dto.ExerciseCreateBase true "Exercises data to create"
// @Success 201 {object} responses.CreatedIDsResponse "Exercises successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/exercise [post]
func (t TrainingHandler) CreateExercises(c *gin.Context) {
	var exercises []dto.ExerciseCreateBase

	if err := c.ShouldBindJSON(&exercises); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	ids, err := t.service.CreateExercises(ctx, t.converter.ExercisesCreateBaseDTOToDomain(exercises))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDsResponse{IDs: ids})
}

// GetExercises
// @Summary Get Exercises with Pagination
// @Description Get exercises with pagination using cursor
// @Tags Trainings
// @Accept json
// @Produce json
// @Param cursor query int false "Cursor for pagination"
// @Param search query string false "Search term"
// @Success 200 {object} dto.ExercisePagination "Return exercises with pagination"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters"
// @Failure 500 "Internal server error"
// @Router /api/training/exercise [get]
func (t TrainingHandler) GetExercises(c *gin.Context) {
	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	search := c.Query("search")

	ctx := c.Request.Context()

	pagination, err := t.service.GetExercises(ctx, search, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pagination)
}

// CreateTrainingBase
// @Summary Create Training Base
// @Description Create a base training
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param training body dto.TrainingCreateBase true "Training data to create"
// @Success 201 {object} responses.CreatedIDsResponse "Training base successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/base [post]
func (t TrainingHandler) CreateTrainingBase(c *gin.Context) {
	var training []dto.TrainingCreateBase

	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	ids, err := t.service.CreateTrainingBases(ctx, t.converter.TrainingCreateBasesDTOToDomain(training))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDsResponse{IDs: ids})
}

// CreateTraining
// @Summary Create Training
// @Description Create a training for a user
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param training body dto.TrainingCreate true "Training data to create"
// @Success 201 {object} responses.CreatedIDIDsResponse "Training successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training [post]
func (t TrainingHandler) CreateTraining(c *gin.Context) {
	var training dto.TrainingCreate

	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	id, ids, err := t.service.CreateTraining(ctx, t.converter.TrainingCreateDTOToDomain(training, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDIDsResponse{
		CreatedIDResponse:  responses.CreatedIDResponse{ID: id},
		CreatedIDsResponse: responses.CreatedIDsResponse{IDs: ids},
	})
}

// SetExerciseStatus
// @Summary Set Exercise Status
// @Description Set the status of an exercise in a training
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param training_id path int true "Training ID"
// @Param exercise_id path int true "Exercise ID"
// @Param status body dto.ExerciseStatusUpdate true "Exercise status to update"
// @Success 200 "Exercise status successfully updated"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/{training_id}/exercise/{exercise_id}/status [patch]
func (t TrainingHandler) SetExerciseStatus(c *gin.Context) {
	trainingIDStr := c.Param("training_id")
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	var statusUpdate dto.ExerciseStatusUpdate
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err = t.service.SetExerciseStatus(ctx, trainingID, exerciseID, statusUpdate.Status)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoExercise):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetTrainings
// @Summary Get Training Covers
// @Description Get training covers with optional search
// @Tags Trainings
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.TrainingCoverPagination "Return training covers"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameters"
// @Failure 500 "Internal server error"
// @Router /api/training [get]
func (t TrainingHandler) GetTrainings(c *gin.Context) {
	search := c.Query("search")
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

	covers, err := t.service.GetTrainingCovers(ctx, search, null.NewInt(0, false), cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, covers)
}

// GetUserTrainings
// @Summary Get User Training Covers
// @Description Get user training covers with optional search and user ID filter
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param search query string false "Search term"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.TrainingCoverPagination "Return training covers"
// @Failure 400 {object} responses.MessageResponse "Bad query or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/user [get]
func (t TrainingHandler) GetUserTrainings(c *gin.Context) {
	search := c.Query("search")
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

	covers, err := t.service.GetTrainingCovers(ctx, search, null.NewInt(int64(userID), true), cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, covers)
}

// GetTraining
// @Summary Get Training
// @Description Get a training by ID
// @Tags Trainings
// @Accept json
// @Produce json
// @Param training_id path int true "Training ID"
// @Success 200 {object} dto.Training "Return training"
// @Failure 400 {object} responses.MessageResponse "Invalid training ID"
// @Failure 404 {object} responses.MessageResponse "No training with such ID"
// @Failure 500 "Internal server error"
// @Router /api/training/{training_id} [get]
func (t TrainingHandler) GetTraining(c *gin.Context) {
	trainingIDStr := c.Param("training_id")
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	training, err := t.service.GetTraining(ctx, trainingID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoTraining):
			c.JSON(http.StatusNotFound, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, training)
}

// GetTrainingsDate
// @Summary Get Trainings Date
// @Description Get trainings by user training IDs
// @Tags Trainings
// @Accept json
// @Produce json
// @Param user_training_ids query []int true "User Training IDs"
// @Success 200 {object} []dto.TrainingDate "Return trainings with dates"
// @Failure 400 {object} responses.MessageResponse "Invalid user training IDs"
// @Failure 404 {object} responses.MessageResponse "No schedule training with such ID"
// @Failure 500 "Internal server error"
// @Router /api/training/date [get]
func (t TrainingHandler) GetTrainingsDate(c *gin.Context) {
	userTrainingIDsStr := c.QueryArray("user_training_ids")
	if len(userTrainingIDsStr) == 0 {
		c.JSON(http.StatusOK, []dto.TrainingDate{})
		return
	}

	var userTrainingIDs []int
	for _, idStr := range userTrainingIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: "Invalid user training ID"})
			return
		}
		userTrainingIDs = append(userTrainingIDs, id)
	}

	ctx := c.Request.Context()

	trainingDates, err := t.service.GetTrainingsDate(ctx, userTrainingIDs)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoTraining):
			c.JSON(http.StatusNotFound, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, trainingDates)
}

// ScheduleTraining
// @Summary Schedule Training
// @Description Schedule training
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param training body dto.ScheduleTraining true "Scheduled training data to create"
// @Success 201 {object} responses.CreatedIDResponse "Scheduled training successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/schedule [post]
func (t TrainingHandler) ScheduleTraining(c *gin.Context) {
	var training dto.ScheduleTraining

	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	id, ids, err := t.service.ScheduleTraining(ctx, t.converter.ScheduleTrainingDTOToDomain(training, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDIDsResponse{
		CreatedIDResponse:  responses.CreatedIDResponse{ID: id},
		CreatedIDsResponse: responses.CreatedIDsResponse{IDs: ids},
	})
}

// GetSchedule
// @Summary Get Schedule
// @Description Get schedule for the specified month
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param month query int true "Month (1-12)"
// @Success 200 {array} dto.Schedule "Return schedule for the month"
// @Failure 400 {object} responses.MessageResponse "Bad month or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/schedule [get]
func (t TrainingHandler) GetSchedule(c *gin.Context) {
	monthStr := c.Query("month")
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	schedule, err := t.service.GetSchedule(ctx, month, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, schedule)
}
