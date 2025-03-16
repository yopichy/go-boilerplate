package testutil

import (
	"testing"
	"webapi/config"
	"webapi/internal/models"

	"github.com/spf13/viper"
)

func SetupTestConfig(t *testing.T) *config.Config {
	cfg, err := config.Load("../../config/testdata")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}
	return cfg
}

func LoadTestWeatherData(t *testing.T) models.WeatherInfo {
	viper.Reset()
	viper.SetConfigName("weather")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../models/testdata")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to load test weather data: %v", err)
	}

	var weather models.WeatherInfo
	weather.Temperature = viper.GetFloat64("weather.temperature")
	weather.Condition = viper.GetString("weather.description")
	weather.Location = viper.GetString("weather.location")
	return weather
}
