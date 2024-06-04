package domain

import "gopkg.in/guregu/null.v3"

type TrainerBase struct {
	FirstName  string
	LastName   string
	Age        int
	Sex        int
	Experience int
	Quote      null.String
}

type TrainerCreate struct {
	TrainerBase
	Email    string
	Password string
}

type TrainerUpdate struct {
	TrainerBase
	ID    int
	Email string
}