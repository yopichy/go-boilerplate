package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webapi/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuth2Authentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.SetupTestConfig(t)

	t.Run("missing token", func(t *testing.T) {
		router := gin.New()
		router.Use(OAuth2Authentication(cfg.OAuth2))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid token in header", func(t *testing.T) {
		router := gin.New()
		router.Use(OAuth2Authentication(cfg.OAuth2))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("valid token in cookie", func(t *testing.T) {
		router := gin.New()
		router.Use(OAuth2Authentication(cfg.OAuth2))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "test-token",
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
