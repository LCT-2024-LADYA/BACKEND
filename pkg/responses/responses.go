package responses

const (
	ResponseBadQuery = "Bad query provided"
	ResponseBadBody  = "Bad body provided"
	ResponseBadPath  = "Bad path parameter provided"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type CreatedIDResponse struct {
	ID int `json:"id"`
}

type CreatedIDsResponse struct {
	IDs []int `json:"ids"`
}
