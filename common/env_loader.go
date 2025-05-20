package common

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnvVariables(logger *zap.Logger) error {
	if err := godotenv.Load(); err != nil {
		logger.Error(
			"Failed to load env variables",
			zap.Error(err),
		)

		return err
	}

	logger.Info(
		"Successfully loaded env variables",
	)

	return nil
}
