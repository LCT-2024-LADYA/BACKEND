package routers

import (
	"BACKEND/internal/delivery/handlers"
	"BACKEND/internal/repository"
	"BACKEND/internal/services"
	"BACKEND/pkg/config"
	"BACKEND/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

func InitRouting(engine *gin.Engine, db *sqlx.DB, jwtUtil utils.JWT, session utils.Session, logger zerolog.Logger) {
	dbResponseTime := time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second

	validate := validator.New()

	// Инициализация репозиториев
	userRepo := repository.InitUserRepo(db)

	// Инициализация сервисов
	userService := services.InitUserService(userRepo, dbResponseTime, logger)
	tokenService := services.InitTokenService(jwtUtil, session)

	// Инициализация хендлеров
	userHandler := handlers.InitAuthHandler(userService, tokenService, validate)

	// Группа маршрутов без middleware
	baseGroup := engine.Group("/api")
	initAuthRouter(baseGroup, userHandler)
}

func initAuthRouter(userGroup *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	authGroup := userGroup.Group("/auth")

	authGroup.POST("/vk", authHandler.AuthorizeVK)
	authGroup.GET("/refresh", authHandler.Refresh)
}
