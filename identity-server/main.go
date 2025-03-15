package main

import (
	"fmt"
	"identity-server/config"
	"identity-server/internal/handlers"
	"identity-server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Connect to default postgres database
	defaultDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Port)

	defaultDb, err := gorm.Open(postgres.Open(defaultDsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to default database: %v", err))
	}

	// Create identity server database if it doesn't exist
	defaultDb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Database))

	// Connect to identity server database
	appDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.Port)

	db, err = gorm.Open(postgres.Open(appDsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to application database: %v", err))
	}

	// Auto migrate schemas
	db.AutoMigrate(&models.User{}, &models.OAuthClient{}, &models.OAuthToken{}, &models.OAuthCode{})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, []byte(cfg.Auth.JWTSecret))
	oauthHandler := handlers.NewOAuthHandler(db)

	r := gin.Default()

	// Public routes (no auth required)
	r.POST("/oauth/clients", oauthHandler.RegisterClient)
	r.GET("/oauth/authorize", oauthHandler.Authorize)

	// Protected routes
	oauth := r.Group("/oauth")
	oauth.Use(authHandler.RequireAuth())
	{
		oauth.POST("/token", oauthHandler.Token)
	}

	// Routes
	r.POST("/users", authHandler.CreateUser)
	r.POST("/login", authHandler.Login)

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
