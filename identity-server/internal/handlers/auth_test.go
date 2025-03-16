package handlers

import (
	"identity-server/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T) *config.Config {
	cfg, err := config.Load("../../../config/testdata")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}
	return cfg
}

func TestAuthHandler_RequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := setupTestConfig(t)

	t.Run("missing auth header", func(t *testing.T) {
		handler := NewAuthHandler(nil, []byte(cfg.Auth.JWTSecret))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		middleware := handler.RequireAuth()
		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
