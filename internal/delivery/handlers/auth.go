package handlers

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/responses"
	"BACKEND/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthHandler struct {
	service      services.Users
	tokenService services.Tokens
	validate     *validator.Validate
}

func InitAuthHandler(
	userService services.Users,
	tokenService services.Tokens,
	validate *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		service:      userService,
		tokenService: tokenService,
		validate:     validate,
	}
}

// AuthorizeVK handles VK authorization callback
// @Summary VK authorization
// @Description Authorize user with VK service
// @Tags Authorization
// @Accept json
// @Produce json
// @Param auth body dto.AuthRequest true "Authorization request body"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad body provided"
// @Failure 401 {object} responses.MessageResponse "Need to auth"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/vk [post]
func (a AuthHandler) AuthorizeVK(c *gin.Context) {
	var user dto.AuthRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := a.validate.Struct(user); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.AuthRequest{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	userID, err := a.service.CreateUserIfNotExistByVK(ctx, user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	tokens, err := a.tokenService.Create(ctx, userID, utils.User)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh handles the refresh token request
// @Summary Refresh tokens
// @Description Refreshes the access and refresh tokens using the provided refresh token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param refresh_token query string true "Refresh token"
// @Success 200 {object} responses.TokenResponse "Return tokens"
// @Failure 400 {object} responses.MessageResponse "Bad query provided"
// @Failure 500 "Internal Server Error"
// @Router /api/auth/refresh [get]
func (a AuthHandler) Refresh(c *gin.Context) {
	refreshToken, ok := c.GetQuery("refresh_token")
	if !ok {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	ctx := c.Request.Context()

	tokens, err := a.tokenService.Refresh(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, errs.NeedToAuth):
			c.JSON(http.StatusUnauthorized, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}
