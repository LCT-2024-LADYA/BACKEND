package main

import (
	"BACKEND/internal/delivery/docs"
	"BACKEND/internal/delivery/routers"
	"BACKEND/pkg/config"
	"BACKEND/pkg/database"
	"BACKEND/pkg/log"
	"BACKEND/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	router.Static("/static", "../static")

	logger, loggerFile := log.InitLoggers()
	defer loggerFile.Close()
	logger.Info().Msg("Logger Initialized")

	config.InitConfig()
	logger.Info().Msg("Config Initialized")

	db := database.GetDB()
	logger.Info().Msg("Database Initialized")

	routers.InitRouting(router, db, utils.InitJWTUtil(), utils.InitRedisSession(), logger)
	logger.Info().Msg("Routing Initialized")

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info().Msg("Swagger Initialized")

	if err := router.Run("0.0.0.0:80"); err != nil {
		panic(fmt.Sprintf("Failed to run client: %s", err.Error()))
	}
}
