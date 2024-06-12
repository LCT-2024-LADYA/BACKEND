package dto

type UserBase struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Age       int    `json:"age" validate:"required,min=14,max=150"`
	Sex       int    `json:"sex" validate:"required,oneof=1 2"`
}

type UserCreate struct {
	UserBase
	Auth
}

type UserUpdate struct {
	UserBase
	Email string `json:"email" validate:"required,email"`
}

type UserCover struct {
	UserBase
	ID       int     `json:"id"`
	PhotoUrl *string `json:"photo_url"`
}

type UserCoverPagination struct {
	Users  []UserCover `json:"objects"`
	Cursor int         `json:"cursor"`
}

type User struct {
	UserCover
	Email string `json:"email"`
}
