package database

import (
	"fmt"
	"webapi/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
}

type database struct {
	db *gorm.DB
}

func getDefaultDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.Port)
}

func getAppDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
}

func New(cfg config.DatabaseConfig) Database {
	// Connect to default postgres database
	defaultDb, err := gorm.Open(postgres.Open(getDefaultDSN(cfg)), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to default database: %v", err))
	}

	// Create application database if it doesn't exist
	defaultDb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database))

	// Connect to application database
	db, err := gorm.Open(postgres.Open(getAppDSN(cfg)), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to application database: %v", err))
	}

	return &database{db: db}
}

func (d *database) GetDB() *gorm.DB {
	return d.db
}
