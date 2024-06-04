package middleware

import (
	"BACKEND/pkg/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserID = "user_id"
)

const (
	AccessToken = "access_token"
)

func (m Middleware) Authorization(userType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(AccessToken)
		if accessToken == "" {
			m.logger.Error().Msg("No jwt provided")
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.MessageResponse{Message: "No jwt provided"})
			return
		}

		userData, isValid, err := m.jwtUtil.Authorize(accessToken, userType)
		if err != nil {
			m.logger.Error().Msg(fmt.Sprintf("Troubles while getting user info from jwt: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.MessageResponse{Message: "Bad jwt provided"})
			return
		}

		if !isValid {
			m.logger.Error().Msg("Access token is expired or invalid or user type mismatch")
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.MessageResponse{Message: "Access token is expired or invalid"})
			return
		}

		c.Set(UserID, userData.ID)
	}
}
