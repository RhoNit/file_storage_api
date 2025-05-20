package main

import (
	"github.com/RhoNit/file_storage_api/common"
	_ "github.com/RhoNit/file_storage_api/docs"
	"github.com/RhoNit/file_storage_api/internal/endpoints"
	"github.com/RhoNit/file_storage_api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/swag"
)

// @title File Storage API
// @version 1.0
// @description APIs for file storage with jwt-authentication

// @contact.name Ranit Biswas
// @contact.email ranitbiswas.cs@gmail.com
//
// @host localhost:8085
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	// intialize zap logger
	zapLogger := common.InitZapLogger()

	// load env vars
	if err := common.LoadEnvVariables(zapLogger); err != nil {
		return
	}

	// initialize the echo engine
	engine := echo.New()

	// echo middlewares
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())

	// initialize handler
	handler := handlers.InitHandler(zapLogger)

	// initiate the routes
	endpoints.InitRoutes(engine, handler)

	// swagger routes
	engine.GET("/swaggger/*", echoSwagger.WrapHandler)
	engine.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(301, "/swagger/index.html")
	})
	// app routes
	// engine.POST("/register", handler.RegisterUserHandler)
	// engine.POST("/login", handler.LoginUserHandler)

	// engine.POST("/upload", handler.UploadFileHandler, common.JWTAuthMiddleware)
	// engine.GET("/storage/remaining", handler.GetRemainingStorageHandler, common.JWTAuthMiddleware)
	// engine.GET("/files", handler.GetUserFilesHandler, common.JWTAuthMiddleware)

	// start the server
	engine.Logger.Fatal(engine.Start("127.0.0.1:8085"))
}
