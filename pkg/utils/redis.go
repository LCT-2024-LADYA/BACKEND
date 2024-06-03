package utils

import (
	"BACKEND/internal/errs"
	"BACKEND/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

type Session interface {
	Set(ctx context.Context, data SessionData) (string, error)
	GetAndUpdate(ctx context.Context, refreshToken string) (string, SessionData, error)
}

type SessionData struct {
	UserID   int    `json:"id"`
	UserType string `json:"type"`
}

type RedisSession struct {
	rdb               *redis.Client
	sessionExpiration time.Duration
	dbResponseTime    time.Duration
}

func InitRedisSession() Session {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString(config.SessionHost),
			viper.GetInt(config.SessionPort),
		),
		Password: viper.GetString(config.SessionPassword),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %s", err.Error()))
	}

	return &RedisSession{
		rdb:               rdb,
		sessionExpiration: time.Duration(viper.GetInt(config.SessionSaveTime)) * time.Hour * 24,
		dbResponseTime:    time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
	}
}

func (r RedisSession) Set(ctx context.Context, data SessionData) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.dbResponseTime)
	defer cancel()

	sessionDataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	key := uuid.New().String()
	err = r.rdb.Set(ctx, key, sessionDataJSON, r.sessionExpiration).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

func (r RedisSession) GetAndUpdate(ctx context.Context, refreshToken string) (string, SessionData, error) {
	var userData SessionData

	ctxGet, cancelGet := context.WithTimeout(ctx, r.dbResponseTime)
	defer cancelGet()

	data, err := r.rdb.Get(ctxGet, refreshToken).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return "", SessionData{}, errs.NeedToAuth
		default:
			return "", SessionData{}, err
		}
	}

	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return "", SessionData{}, err
	}

	key, err := r.Set(ctx, userData)
	if err != nil {
		return "", SessionData{}, err
	}

	ctxDel, cancelDel := context.WithTimeout(ctx, r.dbResponseTime)
	defer cancelDel()

	err = r.rdb.Del(ctxDel, refreshToken).Err()
	if err != nil {
		return "", SessionData{}, err
	}

	return key, userData, nil
}
