package log

import "fmt"

const (
	CreateObject     = "Object `%s` was successfully created with id %d"
	CreateObjects    = "Objects `%s` were successfully created with ids %v"
	GetObjects       = "Objects `%s` were successfully got"
	GetObject        = "Object `%s` with id %d was successfully got"
	UpdateObject     = "Object `%s` with id %d was successfully updated"
	UpdateObjects    = "Objects `%s` with ids %v were successfully updated"
	DeleteObject     = "Object `%s` with id %d was successfully deleted"
	DeleteObjects    = "Object `%s` with ids %v were successfully deleted"
	AuthorizeUser    = "User with email %s authorized"
	AuthorizeTrainer = "Trainer with email %s authorized"
)

const (
	User     = "user"
	Trainer  = "trainer"
	Exercise = "exercise"
	Training = "training"
	Schedule = "schedule"
)

func Normalizer(mainEvent string, args ...any) string {
	return fmt.Sprintf(mainEvent, args...)
}
