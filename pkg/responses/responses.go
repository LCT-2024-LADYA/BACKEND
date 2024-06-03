package responses

const (
	ResponseBadQuery = "Bad query provided"
	ResponseBadBody  = "Bad body provided"
	ResponseBadPath  = "Bad path parameter provided"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type CreatedIDResponse struct {
	ID int `json:"id"`
}

type CreatedIDsResponse struct {
	IDs []int `json:"ids"`
}
