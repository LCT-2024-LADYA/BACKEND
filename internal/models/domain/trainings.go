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

type TrainingCreateBase struct {
	Name        string
	Description string
	Exercises   []int
}

type ExerciseStep struct {
	ID   int
	Step int
}
type TrainingCreate struct {
	UserID      int
	Name        string
	Description string
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

type UserTraining struct {
	Training
	Exercises []Exercise
	Date      time.Time
	TimeStart time.Time
	TimeEnd   time.Time
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

type Schedule struct {
	Date        time.Time
	TrainingIDs []int
}
