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

type ExerciseCreate struct {
	ID     int
	Sets   int
	Reps   int
	Weight int
	Status bool
}

type Exercise struct {
	ExerciseCreateBase
	ExerciseCreate
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

type TrainingCreate struct {
	UserID      int
	Name        string
	Description string
	Exercises   []ExerciseCreate
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
	Exercises []Exercise
}

type TrainingDate struct {
	Training
	Date      time.Time
	TimeStart time.Time
	TimeEnd   time.Time
}

type ScheduleTraining struct {
	TrainingID int
	UserID     int
	Date       time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
}

type Schedule struct {
	Date        time.Time
	TrainingIDs []int
}
