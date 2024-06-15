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
	TrainingCoverTrainerDomainToDTO(training domain.TrainingCoverTrainer) dto.TrainingCoverTrainer
	TrainingCoversTrainerDomainToDTO(trainings []domain.TrainingCoverTrainer) []dto.TrainingCoverTrainer
	TrainingCoverPaginationDomainToDTO(exercise domain.TrainingCoverPagination) dto.TrainingCoverPagination
	TrainingCoverTrainerPaginationDomainToDTO(exercise domain.TrainingCoverTrainerPagination) dto.TrainingCoverTrainerPagination
	TrainingDomainToDTO(training domain.Training) dto.Training
	TrainingTrainerDomainToDTO(training domain.TrainingTrainer) dto.TrainingTrainer
	TrainingDateDomainToDTO(training domain.UserTraining) dto.UserTraining
	TrainingsDateDomainToDTO(trainings []domain.UserTraining) []dto.UserTraining
	TrainingScheduleDomainToDTO(schedule domain.TrainingSchedule) dto.TrainingSchedule
	TrainingSchedulesDomainToDTO(schedules []domain.TrainingSchedule) []dto.TrainingSchedule
	SchedulePlanDomainToDTO(schedule domain.SchedulePlan) dto.SchedulePlan
	SchedulesPlanDomainToDTO(schedule []domain.SchedulePlan) []dto.SchedulePlan
	PlanCoverDomainToDTO(plan domain.PlanCover) dto.PlanCover
	PlanCoversDomainToDTO(plans []domain.PlanCover) []dto.PlanCover
	PlanDomainToDTO(plan domain.Plan) dto.Plan
	ProgressDayDomainToDTO(progress domain.ProgressDay) dto.ProgressDay
	ProgressDaysDomainToDTO(progresses []domain.ProgressDay) []dto.ProgressDay
	ProgressDomainToDTO(progress domain.Progress) dto.Progress
	ProgressesDomainToDTO(progresses []domain.Progress) []dto.Progress
	ProgressPaginationDomainToDTO(progress domain.ProgressPagination) dto.ProgressPagination

	ExerciseCreateBaseDTOToDomain(exercise dto.ExerciseCreateBase) domain.ExerciseCreateBase
	ExercisesCreateBaseDTOToDomain(exercises []dto.ExerciseCreateBase) []domain.ExerciseCreateBase
	TrainingCreateBaseDTOToDomain(training dto.TrainingCreateBase) domain.TrainingCreateBase
	TrainingCreateBasesDTOToDomain(trainings []dto.TrainingCreateBase) []domain.TrainingCreateBase
	TrainingCreateDTOToDomain(training dto.TrainingCreate, userID int) domain.TrainingCreate
	TrainingCreateTrainerDTOToDomain(training dto.TrainingCreateTrainer, trainerID int) domain.TrainingCreateTrainer
	ScheduleTrainingDTOToDomain(training dto.ScheduleTraining, userID int) domain.ScheduleTraining
	ExerciseStepDTOToDomain(exercise dto.ExerciseStep) domain.ExerciseStep
	ExerciseStepsBasesDTOToDomain(exercises []dto.ExerciseStep) []domain.ExerciseStep
	PlanCreateDTOToDomain(plan dto.PlanCreate, userID int) domain.PlanCreate
	SchedulePlanDTOToDomain(plan dto.SchedulePlan) domain.SchedulePlan
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

func (t trainingConverter) TrainingCoverTrainerDomainToDTO(training domain.TrainingCoverTrainer) dto.TrainingCoverTrainer {
	return dto.TrainingCoverTrainer{
		TrainingCover: t.TrainingCoverDomainToDTO(training.TrainingCover),
		WantsPublic:   training.WantsPublic,
		IsConfirm:     training.IsConfirm,
	}
}

func (t trainingConverter) TrainingCoversTrainerDomainToDTO(trainings []domain.TrainingCoverTrainer) []dto.TrainingCoverTrainer {
	result := make([]dto.TrainingCoverTrainer, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingCoverTrainerDomainToDTO(training)
	}

	return result
}

func (t trainingConverter) TrainingCoverPaginationDomainToDTO(exercise domain.TrainingCoverPagination) dto.TrainingCoverPagination {
	return dto.TrainingCoverPagination{
		Trainings: t.TrainingCoversDomainToDTO(exercise.Trainings),
		Cursor:    exercise.Cursor,
	}
}

func (t trainingConverter) TrainingCoverTrainerPaginationDomainToDTO(exercise domain.TrainingCoverTrainerPagination) dto.TrainingCoverTrainerPagination {
	return dto.TrainingCoverTrainerPagination{
		Trainings: t.TrainingCoversTrainerDomainToDTO(exercise.Trainings),
		Cursor:    exercise.Cursor,
	}
}

func (t trainingConverter) TrainingDomainToDTO(training domain.Training) dto.Training {
	return dto.Training{
		TrainingCover: t.TrainingCoverDomainToDTO(training.TrainingCover),
		Exercises:     t.ExercisesBaseStepDomainToDTO(training.Exercises),
	}
}

func (t trainingConverter) TrainingTrainerDomainToDTO(training domain.TrainingTrainer) dto.TrainingTrainer {
	return dto.TrainingTrainer{
		TrainingCoverTrainer: t.TrainingCoverTrainerDomainToDTO(training.TrainingCoverTrainer),
		Exercises:            t.ExercisesBaseStepDomainToDTO(training.Exercises),
	}
}

func (t trainingConverter) TrainingDateDomainToDTO(training domain.UserTraining) dto.UserTraining {
	return dto.UserTraining{
		Training:   t.TrainingDomainToDTO(training.Training),
		Exercises:  t.ExercisesDomainToDTO(training.Exercises),
		TrainingID: training.TrainingID,
		Date:       training.Date,
		TimeStart:  training.TimeStart,
		TimeEnd:    training.TimeEnd,
	}
}

func (t trainingConverter) TrainingsDateDomainToDTO(trainings []domain.UserTraining) []dto.UserTraining {
	result := make([]dto.UserTraining, len(trainings))

	for i, training := range trainings {
		result[i] = t.TrainingDateDomainToDTO(training)
	}

	return result
}

func (t trainingConverter) TrainingScheduleDomainToDTO(schedule domain.TrainingSchedule) dto.TrainingSchedule {
	return dto.TrainingSchedule{
		Date:        schedule.Date,
		TrainingIDs: schedule.TrainingIDs,
	}
}

func (t trainingConverter) TrainingSchedulesDomainToDTO(schedules []domain.TrainingSchedule) []dto.TrainingSchedule {
	result := make([]dto.TrainingSchedule, len(schedules))

	for i, schedule := range schedules {
		result[i] = t.TrainingScheduleDomainToDTO(schedule)
	}

	return result
}

func (t trainingConverter) SchedulePlanDomainToDTO(schedule domain.SchedulePlan) dto.SchedulePlan {
	return dto.SchedulePlan{
		PlanID:    schedule.PlanID,
		DateStart: schedule.DateStart,
		DateEnd:   schedule.DateEnd,
	}
}

func (t trainingConverter) SchedulesPlanDomainToDTO(schedules []domain.SchedulePlan) []dto.SchedulePlan {
	result := make([]dto.SchedulePlan, len(schedules))

	for i, schedule := range schedules {
		result[i] = t.SchedulePlanDomainToDTO(schedule)
	}

	return result
}

func (t trainingConverter) PlanCoverDomainToDTO(plan domain.PlanCover) dto.PlanCover {
	return dto.PlanCover{
		ID:          plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
		Trainings:   plan.Trainings,
	}
}

func (t trainingConverter) PlanCoversDomainToDTO(plans []domain.PlanCover) []dto.PlanCover {
	result := make([]dto.PlanCover, len(plans))

	for i, plan := range plans {
		result[i] = t.PlanCoverDomainToDTO(plan)
	}

	return result
}

func (t trainingConverter) PlanDomainToDTO(plan domain.Plan) dto.Plan {
	return dto.Plan{
		PlanCover: t.PlanCoverDomainToDTO(plan.PlanCover),
		Trainings: t.TrainingCoversDomainToDTO(plan.Trainings),
	}
}

func (t trainingConverter) ProgressDayDomainToDTO(progress domain.ProgressDay) dto.ProgressDay {
	return dto.ProgressDay{
		Sets:   progress.Sets,
		Reps:   progress.Reps,
		Weight: progress.Weight,
		Date:   progress.Date,
	}
}

func (t trainingConverter) ProgressDaysDomainToDTO(progresses []domain.ProgressDay) []dto.ProgressDay {
	result := make([]dto.ProgressDay, len(progresses))

	for i, progress := range progresses {
		result[i] = t.ProgressDayDomainToDTO(progress)
	}

	return result
}

func (t trainingConverter) ProgressDomainToDTO(progress domain.Progress) dto.Progress {
	return dto.Progress{
		Name:       progress.Name,
		Progresses: t.ProgressDaysDomainToDTO(progress.Progresses),
	}
}

func (t trainingConverter) ProgressesDomainToDTO(progresses []domain.Progress) []dto.Progress {
	result := make([]dto.Progress, len(progresses))

	for i, progress := range progresses {
		result[i] = t.ProgressDomainToDTO(progress)
	}

	return result
}

func (t trainingConverter) ProgressPaginationDomainToDTO(progress domain.ProgressPagination) dto.ProgressPagination {
	return dto.ProgressPagination{
		Progresses: t.ProgressesDomainToDTO(progress.Progresses),
		IsMore:     progress.IsMore,
	}
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

func (t trainingConverter) TrainingCreateTrainerDTOToDomain(training dto.TrainingCreateTrainer, trainerID int) domain.TrainingCreateTrainer {
	return domain.TrainingCreateTrainer{
		TrainerID:   trainerID,
		Name:        training.Name,
		Description: training.Description,
		WantsPublic: training.WantsPublic,
		Exercises:   t.ExerciseStepsBasesDTOToDomain(training.Exercises),
	}
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

func (t trainingConverter) PlanCreateDTOToDomain(plan dto.PlanCreate, userID int) domain.PlanCreate {
	return domain.PlanCreate{
		UserID:      userID,
		Name:        plan.Name,
		Description: plan.Description,
		Trainings:   plan.Trainings,
	}
}

func (t trainingConverter) SchedulePlanDTOToDomain(plan dto.SchedulePlan) domain.SchedulePlan {
	return domain.SchedulePlan{
		PlanID:    plan.PlanID,
		DateStart: plan.DateStart,
		DateEnd:   plan.DateEnd,
	}
}
