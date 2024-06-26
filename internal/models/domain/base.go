package domain

type BaseBase struct {
	Name string `json:"name"`
}

type Base struct {
	BaseBase
	ID int `json:"id"`
}

type BaseStatus struct {
	Base
	IsConfirmed bool `json:"is_confirmed"`
}
