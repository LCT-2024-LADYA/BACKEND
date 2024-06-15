package domain

import "time"

type ExerciseCreateBase struct {
	Name             string   `json:"name"`
	Muscle           string   `json:"muscle"`
	AdditionalMuscle string   `json:"additionalMuscle"`
	Type             string   `json:"type"`
	Equipment        string   `json:"equipment"`
	Difficulty       string   `json:"difficulty"`
	Photos           []string `json:"photos"`
}

type ExerciseBase struct {
	ExerciseCreateBase
	ID int
}

type ExerciseBaseStep struct {
	ExerciseBase
	Step int
}

type Exercise struct {
	ExerciseBaseStep
	Sets   int
	Reps   int
	Weight int
	Status bool
}

type ExercisePagination struct {
	Exercises []ExerciseBase
	Cursor    int
}

type ExerciseStep struct {
	ID   int
	Step int
}

type TrainingCreateBase struct {
	Name        string
	Description string
	Exercises   []ExerciseStep
}

type TrainingCreate struct {
	UserID      int
	Name        string
	Description string
	Exercises   []ExerciseStep
}

type TrainingCreateTrainer struct {
	TrainerID   int
	Name        string
	Description string
	WantsPublic bool
	Exercises   []ExerciseStep
}

type TrainingCover struct {
	ID          int
	Name        string
	Description string
	Exercises   int
}

type TrainingCoverPagination struct {
	Trainings []TrainingCover
	Cursor    int
}

type Training struct {
	TrainingCover
	Exercises []ExerciseBaseStep
}

type TrainingCoverTrainer struct {
	TrainingCover
	WantsPublic bool
	IsConfirm   bool
}

type TrainingCoverTrainerPagination struct {
	Trainings []TrainingCoverTrainer
	Cursor    int
}

type TrainingTrainer struct {
	TrainingCoverTrainer
	Exercises []ExerciseBaseStep
}

type UserTraining struct {
	Training
	Exercises  []Exercise
	TrainingID int
	Date       time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
}

type ExerciseDetail struct {
	ExerciseID int
	Sets       int
	Reps       int
	Weight     int
}
type ScheduleTraining struct {
	TrainingID int
	UserID     int
	Date       time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
	Exercises  []ExerciseDetail
}

type TrainingSchedule struct {
	Date        time.Time
	TrainingIDs []int
}

type PlanCreate struct {
	UserID      int
	Name        string
	Description string
	Trainings   []int
}

type PlanCover struct {
	ID          int
	Name        string
	Description string
	Trainings   int
}

type Plan struct {
	PlanCover
	Trainings []TrainingCover
}

type SchedulePlan struct {
	PlanID    int
	DateStart time.Time
	DateEnd   time.Time
}

type ProgressDayString struct {
	Sets   int    `json:"sets"`
	Reps   int    `json:"reps"`
	Weight int    `json:"weight"`
	Date   string `json:"date"`
}

type ProgressDay struct {
	Sets   int
	Reps   int
	Weight int
	Date   time.Time
}

type Progress struct {
	Name       string
	Progresses []ProgressDay
}

type ProgressPagination struct {
	Progresses []Progress
	IsMore     bool
}
