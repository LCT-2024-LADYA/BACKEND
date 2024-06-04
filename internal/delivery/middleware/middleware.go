package middleware

import (
	"BACKEND/pkg/utils"
	"github.com/rs/zerolog"
)

type Middleware struct {
	jwtUtil utils.JWT
	logger  zerolog.Logger
}

func InitMiddleware(
	logger zerolog.Logger,
) *Middleware {
	return &Middleware{
		jwtUtil: utils.InitJWTUtil(),
		logger:  logger,
	}
}
