package handlers

import "go.uber.org/zap"

type Handler struct {
	ZapLogger *zap.Logger
}

func InitHandler(zapLogger *zap.Logger) *Handler {
	return &Handler{
		ZapLogger: zapLogger,
	}
}
