package services

import (
	"BACKEND/pkg/responses"
	"BACKEND/pkg/utils"
	"context"
)

type tokenService struct {
	jwtUtil utils.JWT
	session utils.Session
}

func InitTokenService(jwtUtil utils.JWT, session utils.Session) Tokens {
	return &tokenService{
		jwtUtil: jwtUtil,
		session: session,
	}
}

func (t tokenService) Create(ctx context.Context, userID int, userType string) (responses.TokenResponse, error) {
	accessToken := t.jwtUtil.CreateToken(userID, userType)

	refreshToken, err := t.session.Set(ctx, utils.SessionData{
		UserID:   userID,
		UserType: userType,
	})
	if err != nil {
		return responses.TokenResponse{}, err
	}

	return responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t tokenService) Refresh(ctx context.Context, refreshToken string) (responses.TokenResponse, error) {
	newRefreshToken, data, err := t.session.GetAndUpdate(ctx, refreshToken)
	if err != nil {
		return responses.TokenResponse{}, err
	}

	accessToken := t.jwtUtil.CreateToken(data.UserID, data.UserType)

	return responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
