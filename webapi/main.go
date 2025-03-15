package main

import (
	"fmt"
	"webapi/config"
	_ "webapi/docs"
	"webapi/internal/handlers"
	"webapi/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Web API Boilerplate
// @version         1.0
// @description     A Web API boilerplate with authentication and multitenancy
// @host            localhost:8001
// @BasePath        /
// @securityDefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl http://localhost:8000/oauth/authorize
// @tokenUrl http://localhost:8000/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	healthHandler := handlers.NewHealthHandler()
	weatherHandler := handlers.NewWeatherHandler()
	authHandler := handlers.NewAuthHandler(cfg.OAuth2)

	// Public endpoints (no auth required)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/login", authHandler.Login)
	r.GET("/callback", authHandler.Callback)
	r.GET("/health", healthHandler.HealthCheck)

	// Protected API group (auth required)
	api := r.Group("/api")
	api.Use(middleware.OAuth2Authentication(cfg.OAuth2)) // Use OAuth2 authentication middleware
	{
		api.GET("/weather", weatherHandler.GetWeather)
	}

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
