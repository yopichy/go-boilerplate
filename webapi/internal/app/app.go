package app

import (
	"webapi/config"
	"webapi/internal/database"
	"webapi/internal/handlers"
	"webapi/internal/logger"
	"webapi/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	config *config.Config
	router *gin.Engine
	db     database.Database
	logger logger.Logger
}

func New(cfg *config.Config) *App {
	app := &App{
		config: cfg,
		router: gin.Default(),
	}

	app.setupLogger()
	app.setupDatabase()
	app.setupMiddleware()
	app.setupRoutes()

	return app
}

func (a *App) setupLogger() {
	a.logger = logger.New(a.config.Logging)
}

func (a *App) setupDatabase() {
	a.db = database.New(a.config.Database)
}

func (a *App) setupMiddleware() {
	a.router.Use(middleware.AuthMiddleware())
	a.router.Use(middleware.Tenancy())
	a.router.Use(middleware.LocalizationMiddleware())
	a.router.Use(middleware.ErrorHandler(a.logger))
}

func (a *App) setupRoutes() {
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	weatherHandler := handlers.NewWeatherHandler()
	apiGroup := a.router.Group("/api")
	{
		apiGroup.GET("/weather", weatherHandler.GetWeather)
	}
}

func (a *App) Start() error {
	return a.router.Run()
}
