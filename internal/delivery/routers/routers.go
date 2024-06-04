package routers

import (
	"BACKEND/internal/delivery/handlers"
	"BACKEND/internal/delivery/middleware"
	"BACKEND/internal/repository"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/config"
	"BACKEND/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

func InitRouting(engine *gin.Engine, db *sqlx.DB, middleWarrior *middleware.Middleware, jwtUtil utils.JWT, session utils.Session, logger zerolog.Logger) {
	dbResponseTime := time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second

	validate := validator.New()
	validate.RegisterValidation("password", validators.ValidatePassword)

	// Инициализация репозиториев
	userRepo := repository.InitUserRepo(db)
	trainerRepo := repository.InitTrainerRepo(db)

	// Инициализация сервисов
	userService := services.InitUserService(userRepo, dbResponseTime, logger)
	trainerService := services.InitTrainerService(trainerRepo, dbResponseTime, logger)
	tokenService := services.InitTokenService(jwtUtil, session)

	// Инициализация хендлеров
	authHandler := handlers.InitAuthHandler(userService, trainerService, tokenService, validate)

	// Инициализация middleware
	// userMiddleware := middleWarrior.Authorization(utils.User)
	// trainerMiddleware := middleWarrior.Authorization(utils.Trainer)
	adminMiddleware := middleWarrior.Authorization(utils.Admin)

	// Группа маршрутов без middleware
	baseGroup := engine.Group("/api")
	initAuthRouter(baseGroup, authHandler, adminMiddleware)
}

func initAuthRouter(userGroup *gin.RouterGroup, authHandler *handlers.AuthHandler, adminMiddleware gin.HandlerFunc) {
	authGroup := userGroup.Group("/auth")

	authGroup.POST("register/user", authHandler.RegisterUser)
	authGroup.POST("login/user", authHandler.AuthorizeUser)
	authGroup.POST("register/trainer", adminMiddleware, authHandler.RegisterTrainer)
	authGroup.POST("login/trainer", authHandler.AuthorizeTrainer)
	authGroup.POST("login/admin", authHandler.AuthorizeAdmin)
	authGroup.GET("refresh", authHandler.Refresh)
}
