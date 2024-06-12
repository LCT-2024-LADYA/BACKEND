package domain

import "gopkg.in/guregu/null.v3"

type UserBase struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
	Sex       int    `db:"sex"`
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
	ID       int         `db:"id"`
	PhotoUrl null.String `db:"photo_url"`
}

type UserCoverPagination struct {
	Users  []UserCover
	Cursor int
}

type User struct {
	UserCover
	Email string
}
