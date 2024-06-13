package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type TrainingConverter interface {
	ExerciseCreateBaseDomainToDTO(training domain.ExerciseCreateBase) dto.ExerciseCreateBase
	ExerciseBaseDomainToDTO(exercise domain.ExerciseBase) dto.ExerciseBase
	ExercisesBaseDomainToDTO(exercises []domain.ExerciseBase) []dto.ExerciseBase
	ExerciseBaseStepDomainToDTO(exercise domain.ExerciseBaseStep) dto.ExerciseBaseStep
	ExercisesBaseStepDomainToDTO(exercises []domain.ExerciseBaseStep) []dto.ExerciseBaseStep
	ExerciseDomainToDTO(training domain.Exercise) dto.Exercise
	ExercisesDomainToDTO(exercises []domain.Exercise) []dto.Exercise
	ExercisePaginationDomainToDTO(exercise domain.ExercisePagination) dto.ExercisePagination
	TrainingCoverDomainToDTO(trainer domain.TrainingCover) dto.TrainingCover
	TrainingCoversDomainToDTO(trainings []domain.TrainingCover) []dto.TrainingCover
	TrainingCoverPaginationDomainToDTO(exercise domain.TrainingCoverPagination) dto.TrainingCoverPagination
	TrainingDomainToDTO(training domain.Training) dto.Training
	TrainingDateDomainToDTO(training domain.UserTraining) dto.UserTraining
	TrainingsDateDomainToDTO(trainings []domain.UserTraining) []dto.UserTraining
	ScheduleDomainToDTO(schedule domain.Schedule) dto.Schedule
	SchedulesDomainToDTO(schedules []domain.Schedule) []dto.Schedule

	ExerciseCreateBaseDTOToDomain(exercise dto.ExerciseCreateBase) domain.ExerciseCreateBase
	ExercisesCreateBaseDTOToDomain(exercises []dto.ExerciseCreateBase) []domain.ExerciseCreateBase
	TrainingCreateBaseDTOToDomain(training dto.TrainingCreateBase) domain.TrainingCreateBase
	TrainingCreateBasesDTOToDomain(trainings []dto.TrainingCreateBase) []domain.TrainingCreateBase
	TrainingCreateDTOToDomain(training dto.TrainingCreate, userID int) domain.TrainingCreate
	ScheduleTrainingDTOToDomain(training dto.ScheduleTraining, userID int) domain.ScheduleTraining
	ExerciseStepDTOToDomain(exercise dto.ExerciseStep) domain.ExerciseStep
	ExerciseStepsBasesDTOToDomain(exercises []dto.ExerciseStep) []domain.ExerciseStep
}

type trainingConverter struct{}

func InitTrainingConverter() TrainingConverter {
	return &trainingConverter{}
}

// Domain -> DTO

func (t trainingConverter) ExerciseCreateBaseDomainToDTO(exercise domain.ExerciseCreateBase) dto.ExerciseCreateBase {
	return dto.ExerciseCreateBase{
		Name:             exercise.Name,
		Muscle:           exercise.Muscle,
		AdditionalMuscle: exercise.AdditionalMuscle,
		Type:             exercise.Type,
		Equipment:        exercise.Equipment,
		Difficulty:       exercise.Difficulty,
		Photos:           exercise.Photos,
	}
}

func (t trainingConverter) ExerciseBaseDomainToDTO(exercise domain.ExerciseBase) dto.ExerciseBase {
	return dto.ExerciseBase{
		ExerciseCreateBase: t.ExerciseCreateBaseDomainToDTO(exercise.ExerciseCreateBase),
		ID:                 exercise.ID,
	}
}

func (t trainingConverter) ExercisesBaseDomainToDTO(exercises []domain.ExerciseBase) []dto.ExerciseBase {
	result := make([]dto.ExerciseBase, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseBaseDomainToDTO(exercise)
	}

	return result
}

func (t trainingConverter) ExerciseBaseStepDomainToDTO(exercise domain.ExerciseBaseStep) dto.ExerciseBaseStep {
	return dto.ExerciseBaseStep{
		ExerciseBase: t.ExerciseBaseDomainToDTO(exercise.ExerciseBase),
		Step:         exercise.Step,
	}
}

func (t trainingConverter) ExercisesBaseStepDomainToDTO(exercises []domain.ExerciseBaseStep) []dto.ExerciseBaseStep {
	result := make([]dto.ExerciseBaseStep, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseBaseStepDomainToDTO(exercise)
	}

	return result
}

func (t trainingConverter) ExerciseDomainToDTO(exercise domain.Exercise) dto.Exercise {
	return dto.Exercise{
		ExerciseBaseStep: t.ExerciseBaseStepDomainToDTO(exercise.ExerciseBaseStep),
		Sets:             exercise.Sets,
		Reps:             exercise.Reps,
		Weight:           exercise.Weight,
		Status:           exercise.Status,
	}
}

func (t trainingConverter) ExercisesDomainToDTO(exercises []domain.Exercise) []dto.Exercise {
	result := make([]dto.Exercise, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseDomainToDTO(exercise)
	}

	return result
}

func (t trainingConverter) ExercisePaginationDomainToDTO(exercise domain.ExercisePagination) dto.ExercisePagination {
	return dto.ExercisePagination{
		Exercises: t.ExercisesBaseDomainToDTO(exercise.Exercises),
		Cursor:    exercise.Cursor,
	}
}

func (t trainingConverter) TrainingCoverDomainToDTO(training domain.TrainingCover) dto.TrainingCover {
	return dto.TrainingCover{
		ID:          training.ID,
		Name:        training.Name,
		Description: training.Description,
		Exercises:   training.Exercises,
	}
}

func (t trainingConverter) TrainingCoversDomainToDTO(trainings []domain.TrainingCover) []dto.TrainingCover {
	result := make([]dto.TrainingCover, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingCoverDomainToDTO(training)
	}

	return result
}

func (t trainingConverter) TrainingCoverPaginationDomainToDTO(exercise domain.TrainingCoverPagination) dto.TrainingCoverPagination {
	return dto.TrainingCoverPagination{
		Trainings: t.TrainingCoversDomainToDTO(exercise.Trainings),
		Cursor:    exercise.Cursor,
	}
}

func (t trainingConverter) TrainingDomainToDTO(training domain.Training) dto.Training {
	return dto.Training{
		TrainingCover: t.TrainingCoverDomainToDTO(training.TrainingCover),
		Exercises:     t.ExercisesBaseStepDomainToDTO(training.Exercises),
	}
}

func (t trainingConverter) TrainingDateDomainToDTO(training domain.UserTraining) dto.UserTraining {
	return dto.UserTraining{
		Training:  t.TrainingDomainToDTO(training.Training),
		Exercises: t.ExercisesDomainToDTO(training.Exercises),
		Date:      training.Date,
		TimeStart: training.TimeStart,
		TimeEnd:   training.TimeEnd,
	}
}

func (t trainingConverter) TrainingsDateDomainToDTO(trainings []domain.UserTraining) []dto.UserTraining {
	result := make([]dto.UserTraining, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingDateDomainToDTO(training)
	}

	return result
}

func (t trainingConverter) ScheduleDomainToDTO(schedule domain.Schedule) dto.Schedule {
	return dto.Schedule{
		Date:        schedule.Date,
		TrainingIDs: schedule.TrainingIDs,
	}
}

func (t trainingConverter) SchedulesDomainToDTO(schedules []domain.Schedule) []dto.Schedule {
	result := make([]dto.Schedule, len(schedules))

	for i, schedule := range schedules {
		result[i] = t.ScheduleDomainToDTO(schedule)
	}

	return result
}

// DTO -> Domain

func (t trainingConverter) ExerciseCreateBaseDTOToDomain(exercise dto.ExerciseCreateBase) domain.ExerciseCreateBase {
	return domain.ExerciseCreateBase{
		Name:             exercise.Name,
		Muscle:           exercise.Muscle,
		AdditionalMuscle: exercise.AdditionalMuscle,
		Type:             exercise.Type,
		Equipment:        exercise.Equipment,
		Difficulty:       exercise.Difficulty,
		Photos:           exercise.Photos,
	}
}

func (t trainingConverter) ExercisesCreateBaseDTOToDomain(exercises []dto.ExerciseCreateBase) []domain.ExerciseCreateBase {
	result := make([]domain.ExerciseCreateBase, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseCreateBaseDTOToDomain(exercise)
	}

	return result
}

func (t trainingConverter) TrainingCreateBaseDTOToDomain(training dto.TrainingCreateBase) domain.TrainingCreateBase {
	return domain.TrainingCreateBase{
		Name:        training.Name,
		Description: training.Description,
		Exercises:   t.ExerciseStepsBasesDTOToDomain(training.Exercises),
	}
}

func (t trainingConverter) TrainingCreateBasesDTOToDomain(trainings []dto.TrainingCreateBase) []domain.TrainingCreateBase {
	result := make([]domain.TrainingCreateBase, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingCreateBaseDTOToDomain(training)
	}

	return result
}

func (t trainingConverter) ExerciseStepDTOToDomain(exercise dto.ExerciseStep) domain.ExerciseStep {
	return domain.ExerciseStep{
		ID:   exercise.ID,
		Step: exercise.Step,
	}
}

func (t trainingConverter) ExerciseStepsBasesDTOToDomain(exercises []dto.ExerciseStep) []domain.ExerciseStep {
	result := make([]domain.ExerciseStep, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseStepDTOToDomain(exercise)
	}

	return result
}

func (t trainingConverter) TrainingCreateDTOToDomain(training dto.TrainingCreate, userID int) domain.TrainingCreate {
	return domain.TrainingCreate{
		UserID:      userID,
		Name:        training.Name,
		Description: training.Description,
		Exercises:   t.ExerciseStepsBasesDTOToDomain(training.Exercises),
	}
}

func (t trainingConverter) ExerciseDetailDTOToDomain(detail dto.ExerciseDetail) domain.ExerciseDetail {
	return domain.ExerciseDetail{
		ExerciseID: detail.ExerciseID,
		Sets:       detail.Sets,
		Reps:       detail.Reps,
		Weight:     detail.Weight,
	}
}

func (t trainingConverter) ExercisesDetailBasesDTOToDomain(details []dto.ExerciseDetail) []domain.ExerciseDetail {
	result := make([]domain.ExerciseDetail, len(details))

	for i, detail := range details {
		result[i] = t.ExerciseDetailDTOToDomain(detail)
	}

	return result
}

func (t trainingConverter) ScheduleTrainingDTOToDomain(training dto.ScheduleTraining, userID int) domain.ScheduleTraining {
	return domain.ScheduleTraining{
		UserID:     userID,
		TrainingID: training.TrainingID,
		Date:       training.Date,
		TimeStart:  training.TimeStart,
		TimeEnd:    training.TimeEnd,
		Exercises:  t.ExercisesDetailBasesDTOToDomain(training.Exercises),
	}
}
