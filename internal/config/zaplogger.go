package config

import (
	"os"

	"go.uber.org/zap"
)

func NewZapLogger() (*zap.Logger, error) {
	log, err := zap.NewProduction()
	if os.Getenv("APP_ENV") == "development" {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	defer log.Sync()
	return log, nil
}
