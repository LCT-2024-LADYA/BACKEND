package domain

import "gopkg.in/guregu/null.v3"

type UserBase struct {
	FirstName string
	LastName  string
	Age       int
	Sex       int
}

type UserCreate struct {
	UserBase
	Email    string
	Password string
}

type UserUpdate struct {
	UserBase
	ID    int
	Email string
}

type UserCover struct {
	UserBase
	ID       int
	PhotoUrl null.String
}

type User struct {
	UserCover
	Email string
}
