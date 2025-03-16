package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapi/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWeatherHandler_GetWeather(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful weather fetch", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/weather", nil)

		handler := NewWeatherHandler()
		handler.GetWeather(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.WeatherResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotZero(t, response.Temperature)
		assert.NotEmpty(t, response.Description)
	})
}
