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
	userHandler := handlers.InitUserHandler(userService, validate)
	trainerHandler := handlers.InitTrainerHandler(trainerService, validate)

	// Инициализация middleware
	userMiddleware := middleWarrior.Authorization(utils.User)
	trainerMiddleware := middleWarrior.Authorization(utils.Trainer)
	adminMiddleware := middleWarrior.Authorization(utils.Admin)

	// Группа маршрутов без middleware
	baseGroup := engine.Group("/api")
	initAuthRouter(baseGroup, authHandler, adminMiddleware)
	initUserRouter(baseGroup, userHandler, userMiddleware)
	initTrainerRouter(baseGroup, trainerHandler, trainerMiddleware, adminMiddleware)
}

func initAuthRouter(group *gin.RouterGroup, authHandler *handlers.AuthHandler, adminMiddleware gin.HandlerFunc) {
	authGroup := group.Group("/auth")

	authGroup.POST("register/user", authHandler.RegisterUser)
	authGroup.POST("login/user", authHandler.AuthorizeUser)
	authGroup.POST("register/trainer", adminMiddleware, authHandler.RegisterTrainer)
	authGroup.POST("login/trainer", authHandler.AuthorizeTrainer)
	authGroup.POST("login/admin", authHandler.AuthorizeAdmin)
	authGroup.GET("refresh", authHandler.Refresh)
}

func initUserRouter(group *gin.RouterGroup, userHandler *handlers.UserHandler, userMiddleware gin.HandlerFunc) {
	userGroup := group.Group("/user")

	userGroup.GET("me", userMiddleware, userHandler.Me)
	userGroup.GET(":user_id", userHandler.GetProfile)
	userGroup.PUT("main", userMiddleware, userHandler.UpdateMain)
	userGroup.PUT("photo", userMiddleware, userHandler.UpdatePhoto)
}

func initTrainerRouter(group *gin.RouterGroup, trainerHandler *handlers.TrainerHandler, trainerMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	userGroup := group.Group("/trainer")

	userGroup.GET("me", trainerMiddleware, trainerHandler.Me)
	userGroup.GET(":trainer_id", trainerHandler.GetProfile)
	userGroup.PUT("main", trainerMiddleware, trainerHandler.UpdateMain)
	userGroup.PUT("photo", trainerMiddleware, trainerHandler.UpdatePhoto)
	userGroup.PUT("roles", trainerMiddleware, trainerHandler.UpdateRoles)
	userGroup.PUT("specializations", trainerMiddleware, trainerHandler.UpdateSpecializations)
	userGroup.POST("service", trainerMiddleware, trainerHandler.CreateService)
	userGroup.PUT("service", trainerMiddleware, trainerHandler.UpdateService)
	userGroup.DELETE("service/:service_id", trainerMiddleware, trainerHandler.DeleteService)
	userGroup.POST("achievement", trainerMiddleware, trainerHandler.CreateAchievement)
	userGroup.PUT("achievement/:achievement_id/status", adminMiddleware, trainerHandler.UpdateAchievementStatus)
	userGroup.DELETE("achievement/:achievement_id", trainerMiddleware, trainerHandler.DeleteAchievement)
}
