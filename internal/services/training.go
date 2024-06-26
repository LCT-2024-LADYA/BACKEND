package services

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"context"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type trainingService struct {
	trainingRepo   repository.Trainings
	converter      converters.TrainingConverter
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitTrainingService(
	trainingRepo repository.Trainings,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Trainings {
	return &trainingService{
		trainingRepo:   trainingRepo,
		converter:      converters.InitTrainingConverter(),
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (t trainingService) CreateExercises(ctx context.Context, exercises []domain.ExerciseCreateBase) ([]int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	// Изменение пути к фотографиям
	for i := range exercises {
		for j := range exercises[i].Photos {
			exercises[i].Photos[j] = strings.Replace(exercises[i].Photos[j], "image/", "/static/img/exercises/", 1)
		}
	}

	ids, err := t.trainingRepo.CreateExercises(ctx, exercises)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return []int{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObjects, log.Exercise, ids))

	return ids, nil
}

func (t trainingService) GetExercises(ctx context.Context, search string, cursor int) (dto.ExercisePagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	exercises, err := t.trainingRepo.GetExercises(ctx, search, cursor)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.ExercisePagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Exercise))

	return t.converter.ExercisePaginationDomainToDTO(exercises), nil
}

func (t trainingService) CreateTrainingBases(ctx context.Context, trainings []domain.TrainingCreateBase) ([]int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	id, err := t.trainingRepo.CreateTrainingBases(ctx, trainings)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return []int{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Training, id))

	return id, nil
}

func (t trainingService) CreateTraining(ctx context.Context, training domain.TrainingCreate) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	id, err := t.trainingRepo.CreateTraining(ctx, training)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Training, id))

	return id, nil
}

func (t trainingService) CreateTrainingTrainer(ctx context.Context, training domain.TrainingCreateTrainer) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	id, err := t.trainingRepo.CreateTrainingTrainer(ctx, training)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Training, id))

	return id, nil
}

func (t trainingService) SetExerciseStatus(ctx context.Context, usersTrainingsID, usersExercisesID int, status bool) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainingRepo.SetExerciseStatus(ctx, usersTrainingsID, usersExercisesID, status)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.UpdateObject, log.Training, usersTrainingsID))

	return nil
}

func (t trainingService) GetTrainingCovers(ctx context.Context, search string, cursor int) (dto.TrainingCoverPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	covers, err := t.trainingRepo.GetTrainingCovers(ctx, search, cursor)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.TrainingCoverPagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Training))

	return t.converter.TrainingCoverPaginationDomainToDTO(covers), nil
}

func (t trainingService) GetTrainingCoversByUserID(ctx context.Context, search string, userID, cursor int) (dto.TrainingCoverPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	pagination, err := t.trainingRepo.GetTrainingCoversByUserID(ctx, search, userID, cursor)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.TrainingCoverPagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Training))

	return t.converter.TrainingCoverPaginationDomainToDTO(pagination), nil
}

func (t trainingService) GetTrainingCoversByTrainerID(ctx context.Context, search string, trainerID, cursor int) (dto.TrainingCoverTrainerPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	pagination, err := t.trainingRepo.GetTrainingCoversByTrainerID(ctx, search, trainerID, cursor)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.TrainingCoverTrainerPagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Training))

	return t.converter.TrainingCoverTrainerPaginationDomainToDTO(pagination), nil
}

func (t trainingService) GetTraining(ctx context.Context, trainingID int) (dto.Training, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	training, err := t.trainingRepo.GetTraining(ctx, trainingID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.Training{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObject, log.Training, trainingID))

	return t.converter.TrainingDomainToDTO(training), nil
}

func (t trainingService) GetTrainingTrainer(ctx context.Context, trainingID int) (dto.TrainingTrainer, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	training, err := t.trainingRepo.GetTrainingTrainer(ctx, trainingID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.TrainingTrainer{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObject, log.Training, trainingID))

	return t.converter.TrainingTrainerDomainToDTO(training), nil
}

func (t trainingService) GetScheduleTrainings(ctx context.Context, userTrainingIDs []int) ([]dto.UserTraining, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	training, err := t.trainingRepo.GetScheduleTrainings(ctx, userTrainingIDs)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return []dto.UserTraining{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Training, userTrainingIDs))

	return t.converter.TrainingsDateDomainToDTO(training), nil
}

func (t trainingService) ScheduleTraining(ctx context.Context, training domain.ScheduleTraining) (int, []int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	createdID, createdIDs, err := t.trainingRepo.ScheduleTraining(ctx, training)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, []int{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Schedule, createdID))

	return createdID, createdIDs, nil
}

func (t trainingService) GetSchedule(ctx context.Context, month, userID int) ([]dto.TrainingSchedule, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	schedules, err := t.trainingRepo.GetSchedule(ctx, month, userID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return []dto.TrainingSchedule{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Schedule))

	return t.converter.TrainingSchedulesDomainToDTO(schedules), nil
}

func (t trainingService) DeleteUserTraining(ctx context.Context, trainingID int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainingRepo.DeleteUserTraining(ctx, trainingID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.DeleteObject, log.Training))

	return nil
}

func (t trainingService) DeleteScheduledTraining(ctx context.Context, userTrainingID int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainingRepo.DeleteScheduledTraining(ctx, userTrainingID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.DeleteObject, log.Training, userTrainingID))

	return nil
}

func (t trainingService) CreatePlan(ctx context.Context, plan domain.PlanCreate) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	createdID, err := t.trainingRepo.CreatePlan(ctx, plan)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return 0, err
	}

	t.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Plan, createdID))

	return createdID, nil
}

func (t trainingService) GetPlanCoversByUserID(ctx context.Context, userID int) ([]dto.PlanCover, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	plans, err := t.trainingRepo.GetPlanCoversByUserID(ctx, userID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return []dto.PlanCover{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Plan))

	return t.converter.PlanCoversDomainToDTO(plans), nil
}

func (t trainingService) GetPlan(ctx context.Context, planID int) (dto.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	plan, err := t.trainingRepo.GetPlan(ctx, planID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.Plan{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObject, log.Plan, planID))

	return t.converter.PlanDomainToDTO(plan), nil
}

func (t trainingService) DeletePlan(ctx context.Context, planID int) error {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	err := t.trainingRepo.DeletePlan(ctx, planID)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return err
	}

	t.logger.Info().Msg(log.Normalizer(log.DeleteObject, log.Plan, planID))

	return nil
}

func (t trainingService) GetProgress(ctx context.Context, filters domain.FiltersProgress) (dto.ProgressPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.dbResponseTime)
	defer cancel()

	progress, err := t.trainingRepo.GetProgress(ctx, filters)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return dto.ProgressPagination{}, err
	}

	t.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Schedule))

	return t.converter.ProgressPaginationDomainToDTO(progress), nil
}
