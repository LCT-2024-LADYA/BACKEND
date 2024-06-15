package dto

import "time"

type FiltersTrainerCovers struct {
	Search            string `form:"search"`
	Cursor            int    `form:"cursor"`
	RoleIDs           []int  `form:"role_ids"`
	SpecializationIDs []int  `form:"specialization_ids"`
}

type FiltersProgress struct {
	Search    string    `form:"search"`
	DateStart time.Time `form:"date_start"`
	DateEnd   time.Time `form:"date_end"`
	Page      int       `form:"page"`
}
