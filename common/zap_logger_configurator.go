package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger() *zap.Logger {
	// configure the zap logger
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"

	// create the logger
	logger, err := config.Build()
	if err != nil {
		basicConfig := zap.NewDevelopmentConfig()
		basicConfig.OutputPaths = []string{"stderr"}
		logger, _ := basicConfig.Build()

		logger.Error(
			"Failed to create custom logger, rather using basic configuration",
			zap.Error(err),
		)
	}

	return logger
}
