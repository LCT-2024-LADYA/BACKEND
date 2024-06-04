package services

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/pkg/responses"
	"context"
)

type Users interface {
	Register(ctx context.Context, user domain.UserCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
}

type Trainers interface {
	Register(ctx context.Context, trainer domain.TrainerCreate) (int, error)
	Login(ctx context.Context, auth dto.Auth) (int, error)
	UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error
}

type Tokens interface {
	Create(ctx context.Context, userID int, userType string) (responses.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (responses.TokenResponse, error)
}
