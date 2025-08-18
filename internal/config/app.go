package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *zap.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Redis    *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	config = &BootstrapConfig{
		DB:       config.DB,
		App:      config.App,
		Log:      config.Log,
		Validate: validator.New(),
		Config:   viper.New(),
		Redis:    redisClient,
	}
}
