package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Host string
	Port string
}

func NewRedisClient(cfg RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	redisLogger := logger.With(
		zap.String("storage", "Redis"),
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		redisLogger.Error("redis connection error", zap.Error(err))
		return nil, err
	}

	redisLogger.Info("Redis connect", zap.String("pong", pong))

	return redisClient, err
}
