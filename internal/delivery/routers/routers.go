package routers

import (
	"BACKEND/internal/delivery/chat"
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
	entitiesPerRequest := viper.GetInt(config.EntitiesPerRequest)

	validate := validator.New()
	validate.RegisterValidation("password", validators.ValidatePassword)

	// Инициализация репозиториев
	userRepo := repository.InitUserRepo(db)
	trainerRepo := repository.InitTrainerRepo(db)
	specializationRepo := repository.InitBaseRepo(db, repository.SpecializationTable)
	roleRepo := repository.InitBaseRepo(db, repository.RoleTable)
	trainingRepo := repository.InitTrainingRepo(db, entitiesPerRequest)
	chatRepo := repository.InitChatRepo(db, entitiesPerRequest)

	// Инициализация сервисов
	userService := services.InitUserService(userRepo, dbResponseTime, logger)
	trainerService := services.InitTrainerService(trainerRepo, dbResponseTime, logger)
	tokenService := services.InitTokenService(jwtUtil, session)
	specializationService := services.InitBaseService(specializationRepo, dbResponseTime, logger)
	roleService := services.InitBaseService(roleRepo, dbResponseTime, logger)
	trainingService := services.InitTrainingService(trainingRepo, dbResponseTime, logger)
	chatService := services.InitChatService(chatRepo, dbResponseTime, logger)

	// Инициализация хендлеров
	authHandler := handlers.InitAuthHandler(userService, trainerService, tokenService, validate)
	userHandler := handlers.InitUserHandler(userService, validate)
	trainerHandler := handlers.InitTrainerHandler(trainerService, validate)
	specializationHandler := handlers.InitSpecializationHandler(specializationService, validate)
	roleHandler := handlers.InitRoleHandler(roleService, validate)
	trainingHandler := handlers.InitTrainingsHandler(trainingService)
	chatHandler := handlers.InitChatHandler(chatService)
	serviceHandler := handlers.InitServiceHandler(roleService)

	// Инициализация middleware
	userMiddleware := middleWarrior.Authorization(utils.User)
	trainerMiddleware := middleWarrior.Authorization(utils.Trainer)
	adminMiddleware := middleWarrior.Authorization(utils.Admin)

	// Группа маршрутов
	baseGroup := engine.Group("/api")
	initAuthRouter(baseGroup, authHandler, adminMiddleware)
	initUserRouter(baseGroup, userHandler, userMiddleware)
	initTrainerRouter(baseGroup, trainerHandler, trainerMiddleware, adminMiddleware)
	initRolesRouter(baseGroup, roleHandler, adminMiddleware)
	initSpecializationsRouter(baseGroup, specializationHandler, adminMiddleware)
	initTrainingsRouter(baseGroup, trainingHandler, userMiddleware, adminMiddleware)
	initChatRouter(baseGroup, chatHandler, userMiddleware, trainerMiddleware)
	initServiceRouter(baseGroup, serviceHandler)

	wsGroup := engine.Group("/ws")
	chatServer := chat.NewServer(chatService, jwtUtil, logger)
	go chatServer.Listen()
	wsGroup.GET("", chatServer.ChatHandler)
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

func initRolesRouter(group *gin.RouterGroup, roleHandler *handlers.RoleHandler, adminMiddleware gin.HandlerFunc) {
	roleGroup := group.Group("/role")

	roleGroup.POST("", adminMiddleware, roleHandler.CreateRole)
	roleGroup.GET("", roleHandler.GetRoles)
	roleGroup.DELETE("", adminMiddleware, roleHandler.DeleteRoles)
}

func initSpecializationsRouter(group *gin.RouterGroup, specializationHandler *handlers.SpecializationHandler, adminMiddleware gin.HandlerFunc) {
	specializationGroup := group.Group("/specialization")

	specializationGroup.POST("", adminMiddleware, specializationHandler.CreateSpecialization)
	specializationGroup.GET("", specializationHandler.GetSpecializations)
	specializationGroup.DELETE("", adminMiddleware, specializationHandler.DeleteSpecializations)
}

func initServiceRouter(group *gin.RouterGroup, serviceHandler *handlers.ServiceHandler) {
	serviceGroup := group.Group("/service")

	serviceGroup.GET(":id", serviceHandler.GetServiceByID)
}

func initTrainingsRouter(group *gin.RouterGroup, trainingHandler *handlers.TrainingHandler, userMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	trainingGroup := group.Group("/training")

	trainingGroup.POST("exercise", adminMiddleware, trainingHandler.CreateExercises)
	trainingGroup.GET("exercise", trainingHandler.GetExercises)
	trainingGroup.POST("base", adminMiddleware, trainingHandler.CreateTrainingBase)
	trainingGroup.POST("", userMiddleware, trainingHandler.CreateTraining)
	trainingGroup.PATCH(":training_id/exercise/:exercise_id/status", userMiddleware, trainingHandler.SetExerciseStatus)
	trainingGroup.GET("", trainingHandler.GetTrainings)
	trainingGroup.GET("user", userMiddleware, trainingHandler.GetUserTrainings)
	trainingGroup.GET(":training_id", trainingHandler.GetTraining)
	trainingGroup.GET("date", trainingHandler.GetTrainingsDate)
	trainingGroup.POST("schedule", userMiddleware, trainingHandler.ScheduleTraining)
	trainingGroup.GET("schedule", userMiddleware, trainingHandler.GetSchedule)
}

func initChatRouter(group *gin.RouterGroup, chatHandler *handlers.ChatHandler, userMiddleware gin.HandlerFunc, trainerMiddleware gin.HandlerFunc) {
	chatGroup := group.Group("/chat")

	chatGroup.GET("user", userMiddleware, chatHandler.GetUserChats)
	chatGroup.GET("trainer", trainerMiddleware, chatHandler.GetTrainerChats)
	chatGroup.GET("user/:trainer_id", userMiddleware, chatHandler.GetChatMessageUser)
	chatGroup.GET("trainer/:user_id", trainerMiddleware, chatHandler.GetChatMessageTrainer)
}
