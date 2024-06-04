package domain

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
