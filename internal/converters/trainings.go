package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
)

type TrainingConverter interface {
	ExerciseCreateBaseDomainToDTO(training domain.ExerciseCreateBase) dto.ExerciseCreateBase
	ExerciseBaseDomainToDTO(exercise domain.ExerciseBase) dto.ExerciseBase
	ExercisesBaseDomainToDTO(exercises []domain.ExerciseBase) []dto.ExerciseBase
	ExerciseCreateDomainToDTO(training domain.ExerciseCreate) dto.ExerciseCreate
	ExerciseDomainToDTO(training domain.Exercise) dto.Exercise
	ExercisesDomainToDTO(exercises []domain.Exercise) []dto.Exercise
	ExercisePaginationDomainToDTO(exercise domain.ExercisePagination) dto.ExercisePagination
	TrainingCoverDomainToDTO(trainer domain.TrainingCover) dto.TrainingCover
	TrainingCoversDomainToDTO(trainings []domain.TrainingCover) []dto.TrainingCover
	TrainingCoverPaginationDomainToDTO(exercise domain.TrainingCoverPagination) dto.TrainingCoverPagination
	TrainingDomainToDTO(training domain.Training) dto.Training
	TrainingDateDomainToDTO(training domain.TrainingDate) dto.TrainingDate
	TrainingsDateDomainToDTO(trainings []domain.TrainingDate) []dto.TrainingDate
	ScheduleDomainToDTO(schedule domain.Schedule) dto.Schedule
	SchedulesDomainToDTO(schedules []domain.Schedule) []dto.Schedule

	ExerciseCreateBaseDTOToDomain(exercise dto.ExerciseCreateBase) domain.ExerciseCreateBase
	ExercisesCreateBaseDTOToDomain(exercises []dto.ExerciseCreateBase) []domain.ExerciseCreateBase
	ExerciseCreateDTOToDomain(exercise dto.ExerciseCreate) domain.ExerciseCreate
	ExercisesCreateDTOToDomain(exercises []dto.ExerciseCreate) []domain.ExerciseCreate
	TrainingCreateBaseDTOToDomain(training dto.TrainingCreateBase) domain.TrainingCreateBase
	TrainingCreateBasesDTOToDomain(trainings []dto.TrainingCreateBase) []domain.TrainingCreateBase
	TrainingCreateDTOToDomain(training dto.TrainingCreate, userID int) domain.TrainingCreate
	ScheduleTrainingDTOToDomain(training dto.ScheduleTraining, userID int) domain.ScheduleTraining
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

func (t trainingConverter) ExerciseCreateDomainToDTO(exercise domain.ExerciseCreate) dto.ExerciseCreate {
	return dto.ExerciseCreate{
		ID:     exercise.ID,
		Sets:   exercise.Sets,
		Reps:   exercise.Reps,
		Weight: exercise.Weight,
		Status: exercise.Status,
	}
}

func (t trainingConverter) ExerciseDomainToDTO(exercise domain.Exercise) dto.Exercise {
	return dto.Exercise{
		ExerciseCreateBase: t.ExerciseCreateBaseDomainToDTO(exercise.ExerciseCreateBase),
		ExerciseCreate:     t.ExerciseCreateDomainToDTO(exercise.ExerciseCreate),
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
		Exercises:     t.ExercisesDomainToDTO(training.Exercises),
	}
}

func (t trainingConverter) TrainingDateDomainToDTO(training domain.TrainingDate) dto.TrainingDate {
	return dto.TrainingDate{
		Training:  t.TrainingDomainToDTO(training.Training),
		Date:      training.Date,
		TimeStart: training.TimeStart,
		TimeEnd:   training.TimeEnd,
	}
}

func (t trainingConverter) TrainingsDateDomainToDTO(trainings []domain.TrainingDate) []dto.TrainingDate {
	result := make([]dto.TrainingDate, len(trainings))

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

func (t trainingConverter) ExerciseCreateDTOToDomain(exercise dto.ExerciseCreate) domain.ExerciseCreate {
	return domain.ExerciseCreate{
		ID:     exercise.ID,
		Sets:   exercise.Sets,
		Reps:   exercise.Reps,
		Weight: exercise.Weight,
		Status: exercise.Status,
	}
}

func (t trainingConverter) ExercisesCreateDTOToDomain(exercises []dto.ExerciseCreate) []domain.ExerciseCreate {
	result := make([]domain.ExerciseCreate, len(exercises))

	for i, exercise := range exercises {
		result[i] = t.ExerciseCreateDTOToDomain(exercise)
	}

	return result
}

func (t trainingConverter) TrainingCreateBaseDTOToDomain(training dto.TrainingCreateBase) domain.TrainingCreateBase {
	return domain.TrainingCreateBase{
		Name:        training.Name,
		Description: training.Description,
		Exercises:   training.Exercises,
	}
}

func (t trainingConverter) TrainingCreateBasesDTOToDomain(trainings []dto.TrainingCreateBase) []domain.TrainingCreateBase {
	result := make([]domain.TrainingCreateBase, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingCreateBaseDTOToDomain(training)
	}

	return result
}

func (t trainingConverter) TrainingCreateDTOToDomain(training dto.TrainingCreate, userID int) domain.TrainingCreate {
	return domain.TrainingCreate{
		UserID:      userID,
		Name:        training.Name,
		Description: training.Description,
		Exercises:   t.ExercisesCreateDTOToDomain(training.Exercises),
	}
}

func (t trainingConverter) ScheduleTrainingDTOToDomain(training dto.ScheduleTraining, userID int) domain.ScheduleTraining {
	return domain.ScheduleTraining{
		UserID:     userID,
		TrainingID: training.TrainingID,
		Date:       training.Date,
		TimeStart:  training.TimeStart,
		TimeEnd:    training.TimeEnd,
	}
}
