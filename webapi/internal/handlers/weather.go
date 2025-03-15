package handlers

import (
	"net/http"
	"webapi/internal/models"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct{}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{}
}

// GetWeather godoc
// @Summary Get weather information
// @Description Get current weather information
// @Tags weather
// @Accept json
// @Produce json
// @Security OAuth2Implicit
// @Success 200 {object} models.WeatherInfo
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/weather [get]
func (h *WeatherHandler) GetWeather(c *gin.Context) {
	weather := models.WeatherInfo{
		Temperature: 25,
		Condition:   "Sunny",
		Location:    "New York",
	}

	c.JSON(http.StatusOK, weather)
}

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Location    string  `json:"location"`
}
