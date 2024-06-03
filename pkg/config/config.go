package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	DBName         = "DB_NAME"
	DBUser         = "DB_USER"
	DBPassword     = "DB_PASSWORD"
	DBPort         = "DB_PORT"
	DBHost         = "DB_HOST"
	DBResponseTime = "DB_RESPONSE_TIME"

	SessionHost     = "SESSION_HOST"
	SessionPort     = "SESSION_PORT"
	SessionPassword = "SESSION_PASSWORD"
	SessionSaveTime = "SESSION_RESPONSE_TIME"

	JWTExpirationTime = "JWT_EXPIRATION_TIME"
	JWTSecret         = "JWT_SECRET"

	EntitiesPerRequest = "ENTITIES_PER_REQUEST"
)

func InitConfig() {
	envPath, _ := os.Getwd()
	envPath = filepath.Join(envPath, "..") // workdir is cmd

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envPath)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to init config from file. Error:%v", err.Error())
	}
}
