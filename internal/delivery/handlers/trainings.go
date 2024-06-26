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
	"net/http"
	"strconv"
)

type TrainingHandler struct {
	service         services.Trainings
	converter       converters.TrainingConverter
	filterConverter converters.FilterConverter
}

func InitTrainingsHandler(
	service services.Trainings,
) *TrainingHandler {
	return &TrainingHandler{
		service:         service,
		converter:       converters.InitTrainingConverter(),
		filterConverter: converters.InitFilterConverter(),
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
		c.Status(http.StatusInternalServerError)
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
// @Success 201 {object} responses.CreatedIDResponse "Training successfully created"
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

	id, err := t.service.CreateTraining(ctx, t.converter.TrainingCreateDTOToDomain(training, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// CreateTrainingTrainer
// @Summary Create Training Trainer
// @Description Create a training for a trainer
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param training body dto.TrainingCreateTrainer true "Training data to create"
// @Success 201 {object} responses.CreatedIDResponse "Training successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/trainer [post]
func (t TrainingHandler) CreateTrainingTrainer(c *gin.Context) {
	var training dto.TrainingCreateTrainer

	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	trainerID := c.GetInt(middleware.UserID)

	id, err := t.service.CreateTrainingTrainer(ctx, t.converter.TrainingCreateTrainerDTOToDomain(training, trainerID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
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

	covers, err := t.service.GetTrainingCovers(ctx, search, cursor)
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

	covers, err := t.service.GetTrainingCoversByUserID(ctx, search, userID, cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, covers)
}

// GetTrainerTrainings
// @Summary Get Trainer Training Covers
// @Description Get trainer training covers with optional search and trainer ID filter
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param search query string false "Search term"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.TrainingCoverTrainerPagination "Return training covers"
// @Failure 400 {object} responses.MessageResponse "Bad query or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/trainer [get]
func (t TrainingHandler) GetTrainerTrainings(c *gin.Context) {
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

	trainerID := c.GetInt(middleware.UserID)

	covers, err := t.service.GetTrainingCoversByTrainerID(ctx, search, trainerID, cursor)
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

// GetTrainingTrainer
// @Summary Get Training Trainer
// @Description Get a training by ID for trainer
// @Tags Trainings
// @Accept json
// @Produce json
// @Param training_id path int true "Training ID"
// @Success 200 {object} dto.TrainingTrainer "Return training"
// @Failure 400 {object} responses.MessageResponse "Invalid training ID"
// @Failure 404 {object} responses.MessageResponse "No training with such ID"
// @Failure 500 "Internal server error"
// @Router /api/training/{training_id}/trainer [get]
func (t TrainingHandler) GetTrainingTrainer(c *gin.Context) {
	trainingIDStr := c.Param("training_id")
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	training, err := t.service.GetTrainingTrainer(ctx, trainingID)
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

// GetScheduleTrainings
// @Summary Get Trainings Date
// @Description Get trainings by user training IDs
// @Tags Trainings
// @Accept json
// @Produce json
// @Param user_training_ids query []int true "User Training IDs"
// @Success 200 {object} []dto.UserTraining "Return trainings with dates"
// @Failure 400 {object} responses.MessageResponse "Invalid user training IDs"
// @Failure 404 {object} responses.MessageResponse "No schedule training with such ID"
// @Failure 500 "Internal server error"
// @Router /api/training/date [get]
func (t TrainingHandler) GetScheduleTrainings(c *gin.Context) {
	userTrainingIDsStr := c.QueryArray("user_training_ids")
	if len(userTrainingIDsStr) == 0 {
		c.JSON(http.StatusOK, []dto.UserTraining{})
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

	trainingDates, err := t.service.GetScheduleTrainings(ctx, userTrainingIDs)
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
// @Success 200 {object} []dto.TrainingSchedule "Return schedule for the month"
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

// DeleteUserTraining
// @Summary Delete User Training
// @Description Delete a user training
// @Tags Trainings
// @Accept json
// @Produce json
// @Param training_id path int true "Training ID"
// @Success 200 "Training deleted successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid training ID"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/user/{training_id} [delete]
func (t TrainingHandler) DeleteUserTraining(c *gin.Context) {
	trainingIDStr := c.Param("training_id")
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil || trainingID <= 0 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	err = t.service.DeleteUserTraining(ctx, trainingID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteScheduledTraining
// @Summary Delete Scheduled Training
// @Description Delete a scheduled training for a specific date
// @Tags Trainings
// @Accept json
// @Produce json
// @Param user_training_id path int true "User Training ID"
// @Success 200 "Scheduled training deleted successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid user training ID"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/schedule/{user_training_id} [delete]
func (t TrainingHandler) DeleteScheduledTraining(c *gin.Context) {
	userTrainingIDStr := c.Param("user_training_id")
	userTrainingID, err := strconv.Atoi(userTrainingIDStr)
	if err != nil || userTrainingID <= 0 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	err = t.service.DeleteScheduledTraining(ctx, userTrainingID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// CreatePlanUser
// @Summary Create Plan
// @Description Create a new plan
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param plan body dto.PlanCreate true "Plan data to create"
// @Success 201 {object} responses.CreatedIDResponse "Plan successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/plan/user [post]
func (t TrainingHandler) CreatePlanUser(c *gin.Context) {
	var plan dto.PlanCreate

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	id, err := t.service.CreatePlan(ctx, t.converter.PlanCreateDTOToDomain(plan, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// CreatePlanTrainer
// @Summary Create Plan Trainer
// @Description Create a new plan for user from trainer
// @Tags Trainings
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param plan body dto.PlanCreate true "Plan data to create"
// @Success 201 {object} responses.CreatedIDResponse "Plan successfully created"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 500 "Internal server error"
// @Router /api/training/plan/user/{user_id} [post]
func (t TrainingHandler) CreatePlanTrainer(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	var plan dto.PlanCreate

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	id, err := t.service.CreatePlan(ctx, t.converter.PlanCreateDTOToDomain(plan, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: id})
}

// GetPlanCoversByUserID
// @Summary Get Plan Covers by User ID
// @Description Get plan covers by user ID
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Success 200 {object} []dto.PlanCover "Return plan covers"
// @Failure 400 {object} responses.MessageResponse "Bad query or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/plan/user [get]
func (t TrainingHandler) GetPlanCoversByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt(middleware.UserID)

	covers, err := t.service.GetPlanCoversByUserID(ctx, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, covers)
}

// GetPlan
// @Summary Get Plan
// @Description Get a plan by ID
// @Tags Trainings
// @Accept json
// @Produce json
// @Param plan_id path int true "Plan ID"
// @Success 200 {object} dto.Plan "Return plan"
// @Failure 400 {object} responses.MessageResponse "Invalid plan ID"
// @Failure 404 {object} responses.MessageResponse "No plan with such ID"
// @Failure 500 "Internal server error"
// @Router /api/training/plan/{plan_id} [get]
func (t TrainingHandler) GetPlan(c *gin.Context) {
	planIDStr := c.Param("plan_id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	plan, err := t.service.GetPlan(ctx, planID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, plan)
}

// DeletePlan
// @Summary Delete Plan
// @Description Delete a plan by ID
// @Tags Trainings
// @Accept json
// @Produce json
// @Param plan_id path int true "Plan ID"
// @Success 200 "Plan deleted successfully"
// @Failure 400 {object} responses.MessageResponse "Invalid plan ID"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/plan/{plan_id} [delete]
func (t TrainingHandler) DeletePlan(c *gin.Context) {
	planIDStr := c.Param("plan_id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil || planID <= 0 {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	ctx := c.Request.Context()

	err = t.service.DeletePlan(ctx, planID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// GetProgress
// @Summary Get Progress
// @Description Get progress with pagination
// @Tags Trainings
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param search query string false "Search term"
// @Param date_start query string true "Start date"
// @Param date_end query string true "End date"
// @Param page query int false "Page number"
// @Success 200 {object} dto.ProgressPagination "List of progress with pagination"
// @Failure 400 {object} responses.MessageResponse "Bad query or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/training/progress [get]
func (t TrainingHandler) GetProgress(c *gin.Context) {
	var filters dto.FiltersProgress

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	if filters.Page == 0 {
		filters.Page = 1
	}

	userID := c.GetInt(middleware.UserID)

	progressPagination, err := t.service.GetProgress(c.Request.Context(), t.filterConverter.FiltersProgressDTOToDomain(filters, userID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, progressPagination)
}
