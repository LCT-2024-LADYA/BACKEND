package dto

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	Message   *string   `json:"message"`
	Service   *int      `json:"service_id"`
	IsToUser  bool      `json:"is_to_user"`
	Time      time.Time `json:"time"`
}

type MessageCreate struct {
	UserID    int         `json:"user_id"`
	TrainerID int         `json:"trainer_id"`
	Message   null.String `json:"message"`
	ServiceID null.Int    `json:"service_id"`
	IsToUser  bool        `json:"is_to_user"`
}

type MessagePagination struct {
	Messages []Message `json:"objects"`
	Cursor   int       `json:"cursor"`
}

type Chat struct {
	ID              int       `json:"id"`
	PhotoUrl        *string   `json:"photo_url"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	LastMessage     string    `json:"last_message"`
	TimeLastMessage time.Time `json:"time_last_message"`
}
