package dto

type FiltersTrainerCovers struct {
	Search            string `form:"search"`
	Cursor            int    `form:"cursor"`
	RoleIDs           []int  `form:"role_ids"`
	SpecializationIDs []int  `form:"specialization_ids"`
}
