package repository

import (
	"BACKEND/internal/models/domain"
	"context"
)

type Users interface {
	CreateVK(ctx context.Context, user domain.UserCreateVK) (int, error)

	CheckIfExistByVKID(ctx context.Context, VKID int) (int, error)
}
