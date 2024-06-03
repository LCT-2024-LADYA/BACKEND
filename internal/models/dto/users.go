package dto

type AuthRequest struct {
	Token string `json:"access_token" validate:"required"`
	VKID  int    `json:"vk_id" validate:"required,min=1"`
	Email string `json:"email" validate:"required,email"`
}
