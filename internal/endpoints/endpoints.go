package endpoints

import (
	"github.com/RhoNit/file_storage_api/common"
	"github.com/RhoNit/file_storage_api/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, h *handlers.Handler) {
	// app routes
	api := e.Group("/api")
	{
		// public routes
		api.POST("/register", h.RegisterUserHandler)
		api.POST("/login", h.LoginUserHandler)

		// protected routes
		protected := api.Group("")
		protected.Use(common.JWTAuthMiddleware)
		{
			protected.POST("/upload", h.UploadFileHandler)
			protected.GET("/storage/remaining", h.GetRemainingStorageHandler)
			protected.GET("/files", h.GetUserFilesHandler)
		}
	}

	h.ZapLogger.Info(
		"Routes have been initialized successfully",
	)
}
