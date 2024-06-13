package dto

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
	ID int `json:"id"`
}

type ExerciseBaseStep struct {
	ExerciseBase
	Step int
}

type Exercise struct {
	ExerciseBaseStep
	Sets   int  `json:"sets"`
	Reps   int  `json:"reps"`
	Weight int  `json:"weight"`
	Status bool `json:"status"`
}

type ExercisePagination struct {
	Exercises []ExerciseBase `json:"objects"`
	Cursor    int            `json:"cursor"`
}

type ExerciseStep struct {
	ID   int `json:"id"`
	Step int `json:"step"`
}

type TrainingCreateBase struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Exercises   []ExerciseStep `json:"exercises"`
}

type TrainingCreate struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Exercises   []ExerciseStep `json:"exercises"`
}

type TrainingCover struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Exercises   int    `json:"exercises"`
}

type TrainingCoverPagination struct {
	Trainings []TrainingCover `json:"objects"`
	Cursor    int             `json:"cursor"`
}

type Training struct {
	TrainingCover
	Exercises []ExerciseBaseStep `json:"exercises"`
}

type UserTraining struct {
	Training
	Exercises []Exercise
	Date      time.Time `json:"date"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
}

type ExerciseDetail struct {
	ExerciseID int `json:"id"`
	Sets       int `json:"sets"`
	Reps       int `json:"reps"`
	Weight     int `json:"weight"`
}
type ScheduleTraining struct {
	TrainingID int       `json:"id"`
	Date       time.Time `json:"date"`
	TimeStart  time.Time `json:"time_start"`
	TimeEnd    time.Time `json:"time_end"`
	Exercises  []ExerciseDetail
}

type Schedule struct {
	Date        time.Time `json:"date"`
	TrainingIDs []int     `json:"user_training_ids"`
}

type ExerciseStatusUpdate struct {
	Status bool `json:"status"`
}
