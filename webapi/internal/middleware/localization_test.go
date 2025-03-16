package middleware

import (
	"testing"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLocalizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("default language", func(t *testing.T) {
		router := gin.New()
		router.Use(LocalizationMiddleware())

		var capturedLang string
		router.GET("/test", func(c *gin.Context) {
			lang, exists := c.Get("language")
			if exists {
				capturedLang = lang.(string)
			}
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, "en", capturedLang)
	})

	t.Run("custom language header", func(t *testing.T) {
		router := gin.New()
		router.Use(LocalizationMiddleware())

		var capturedLang string
		router.GET("/test", func(c *gin.Context) {
			lang, exists := c.Get("language")
			if exists {
				capturedLang = lang.(string)
			}
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Language", "id-ID")
		router.ServeHTTP(w, req)

		assert.Equal(t, "id", capturedLang)
	})
}
