package domain

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type MessageGet struct {
	To        int     `json:"to"`
	Message   *string `json:"message"`
	ServiceID *int    `json:"service_id"`
}

type Message struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	Message   *string   `json:"message"`
	ServiceID *int      `json:"service_id"`
	IsToUser  bool      `json:"is_to_user"`
	Time      time.Time `json:"time"`
}

type MessageCreate struct {
	UserID    int
	TrainerID int
	Message   null.String
	ServiceID null.Int
	IsToUser  bool
}

type MessagePagination struct {
	Messages []Message
	Cursor   int
}

type Chat struct {
	ID              int
	PhotoUrl        null.String
	FirstName       string
	LastName        string
	LastMessage     string
	TimeLastMessage time.Time
}
