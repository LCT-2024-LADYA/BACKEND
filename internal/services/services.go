package services

import (
	"BACKEND/internal/models/dto"
	"BACKEND/pkg/responses"
	"context"
)

type Users interface {
	CreateUserIfNotExistByVK(ctx context.Context, user dto.AuthRequest) (int, error)
}

type Tokens interface {
	Create(ctx context.Context, userID int, userType string) (responses.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (responses.TokenResponse, error)
}
