package repository

import (
	"BACKEND/internal/models/domain"
	"context"
)

type Users interface {
	Create(ctx context.Context, user domain.UserCreate) (int, error)
	GetSecure(ctx context.Context, email string) (int, string, error)
}

type Trainers interface {
	Create(ctx context.Context, trainer domain.TrainerCreate) (int, error)
	GetSecure(ctx context.Context, email string) (int, string, error)
	UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error
}
