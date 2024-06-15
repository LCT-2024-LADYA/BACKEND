package domain

import "time"

type FiltersTrainerCovers struct {
	Search            string
	Cursor            int
	RoleIDs           []int
	SpecializationIDs []int
}

type FiltersProgress struct {
	UserID    int
	Search    string
	DateStart time.Time
	DateEnd   time.Time
	Page      int
}
