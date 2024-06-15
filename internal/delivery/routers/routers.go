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
	userRepo := repository.InitUserRepo(db, entitiesPerRequest)
	trainerRepo := repository.InitTrainerRepo(db, entitiesPerRequest)
	specializationRepo := repository.InitBaseRepo(db, repository.SpecializationTable)
	roleRepo := repository.InitBaseRepo(db, repository.RoleTable)
	serviceRepo := repository.InitUserTrainerServicesRepo(db, entitiesPerRequest)
	trainingRepo := repository.InitTrainingRepo(db, entitiesPerRequest)
	chatRepo := repository.InitChatRepo(db, entitiesPerRequest)

	// Инициализация сервисов
	userService := services.InitUserService(userRepo, dbResponseTime, logger)
	trainerService := services.InitTrainerService(trainerRepo, dbResponseTime, logger)
	tokenService := services.InitTokenService(jwtUtil, session)
	specializationService := services.InitBaseService(specializationRepo, dbResponseTime, logger)
	roleService := services.InitBaseService(roleRepo, dbResponseTime, logger)
	serviceService := services.InitUsersTrainersServicesService(serviceRepo, dbResponseTime, logger)
	trainingService := services.InitTrainingService(trainingRepo, dbResponseTime, logger)
	chatService := services.InitChatService(chatRepo, dbResponseTime, logger)

	// Инициализация хендлеров
	authHandler := handlers.InitAuthHandler(userService, trainerService, tokenService, validate)
	userHandler := handlers.InitUserHandler(userService, validate)
	trainerHandler := handlers.InitTrainerHandler(trainerService, validate)
	specializationHandler := handlers.InitSpecializationHandler(specializationService, validate)
	roleHandler := handlers.InitRoleHandler(roleService, validate)
	userTrainerServiceHandler := handlers.InitUserTrainerServiceHandler(serviceService)
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
	initUserTrainerServicesRouter(baseGroup, userTrainerServiceHandler, userMiddleware, trainerMiddleware)
	initTrainingsRouter(baseGroup, trainingHandler, userMiddleware, trainerMiddleware, adminMiddleware)
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
	userGroup.GET("", userHandler.GetCovers)
	userGroup.PUT("main", userMiddleware, userHandler.UpdateMain)
	userGroup.PUT("photo", userMiddleware, userHandler.UpdatePhoto)
}

func initTrainerRouter(group *gin.RouterGroup, trainerHandler *handlers.TrainerHandler, trainerMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	userGroup := group.Group("/trainer")

	userGroup.GET("me", trainerMiddleware, trainerHandler.Me)
	userGroup.GET(":trainer_id", trainerHandler.GetProfile)
	userGroup.GET("", trainerHandler.GetCovers)
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

func initUserTrainerServicesRouter(group *gin.RouterGroup, serviceHandler *handlers.UserTrainerServiceHandler, userMiddleware gin.HandlerFunc, trainerMiddleware gin.HandlerFunc) {
	serviceGroup := group.Group("/service")

	serviceGroup.POST("", trainerMiddleware, serviceHandler.CreateService)
	serviceGroup.POST("schedule", serviceHandler.ScheduleService)
	serviceGroup.GET("schedule/:month", trainerMiddleware, serviceHandler.GetSchedule)
	serviceGroup.GET("schedule", serviceHandler.GetSchedulesByIDs)
	serviceGroup.DELETE("schedule/:schedule_id", serviceHandler.DeleteScheduled)
	serviceGroup.GET("trainer", trainerMiddleware, serviceHandler.GetTrainerServices)
	serviceGroup.GET("user", userMiddleware, serviceHandler.GetUserServices)
	serviceGroup.PUT("status/:service_id", serviceHandler.UpdateStatus)
	serviceGroup.DELETE(":service_id", serviceHandler.DeleteService)
}

func initServiceRouter(group *gin.RouterGroup, serviceHandler *handlers.ServiceHandler) {
	serviceGroup := group.Group("/service")

	serviceGroup.GET(":id", serviceHandler.GetServiceByID)
}

func initTrainingsRouter(group *gin.RouterGroup, trainingHandler *handlers.TrainingHandler, userMiddleware, trainerMiddleware, adminMiddleware gin.HandlerFunc) {
	trainingGroup := group.Group("/training")

	trainingGroup.POST("exercise", adminMiddleware, trainingHandler.CreateExercises)
	trainingGroup.GET("exercise", trainingHandler.GetExercises)
	trainingGroup.POST("base", adminMiddleware, trainingHandler.CreateTrainingBase)
	trainingGroup.POST("", userMiddleware, trainingHandler.CreateTraining)
	trainingGroup.POST("trainer", trainerMiddleware, trainingHandler.CreateTrainingTrainer)
	trainingGroup.PATCH(":training_id/exercise/:exercise_id/status", userMiddleware, trainingHandler.SetExerciseStatus)
	trainingGroup.GET("", trainingHandler.GetTrainings)
	trainingGroup.GET("user", userMiddleware, trainingHandler.GetUserTrainings)
	trainingGroup.GET("trainer", trainerMiddleware, trainingHandler.GetTrainerTrainings)
	trainingGroup.GET(":training_id", trainingHandler.GetTraining)
	trainingGroup.GET(":training_id/trainer", trainingHandler.GetTrainingTrainer)
	trainingGroup.GET("date", trainingHandler.GetScheduleTrainings)
	trainingGroup.POST("schedule", userMiddleware, trainingHandler.ScheduleTraining)
	trainingGroup.GET("schedule", userMiddleware, trainingHandler.GetSchedule)
	trainingGroup.DELETE("user/:training_id", trainingHandler.DeleteUserTraining)
	trainingGroup.DELETE("schedule/:user_training_id", trainingHandler.DeleteScheduledTraining)

	trainingGroup.POST("plan/user", userMiddleware, trainingHandler.CreatePlanUser)
	trainingGroup.POST("plan/user/:user_id", trainingHandler.CreatePlanTrainer)
	trainingGroup.GET("plan/user", userMiddleware, trainingHandler.GetPlanCoversByUserID)
	trainingGroup.GET("plan/:plan_id", trainingHandler.GetPlan)
	trainingGroup.DELETE("plan/:plan_id", trainingHandler.DeletePlan)
}

func initChatRouter(group *gin.RouterGroup, chatHandler *handlers.ChatHandler, userMiddleware gin.HandlerFunc, trainerMiddleware gin.HandlerFunc) {
	chatGroup := group.Group("/chat")

	chatGroup.GET("user", userMiddleware, chatHandler.GetUserChats)
	chatGroup.GET("trainer", trainerMiddleware, chatHandler.GetTrainerChats)
	chatGroup.GET("user/:trainer_id", userMiddleware, chatHandler.GetChatMessageUser)
	chatGroup.GET("trainer/:user_id", trainerMiddleware, chatHandler.GetChatMessageTrainer)
}
