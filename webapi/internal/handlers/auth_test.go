package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webapi/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful login redirect", func(t *testing.T) {
		cfg := config.OAuth2Config{
			AuthServerURL: "http://localhost:8000",
			ClientID:      "test-client",
			RedirectURL:   "http://localhost:8080/callback",
		}
		handler := NewAuthHandler(cfg)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/login", nil)

		handler.Login(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "oauth/authorize")

		cookie := w.Result().Cookies()[0]
		assert.Equal(t, "oauth_state", cookie.Name)
		assert.NotEmpty(t, cookie.Value)
	})
}

func TestAuthHandler_Callback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid state", func(t *testing.T) {
		cfg := config.OAuth2Config{
			AuthServerURL: "http://localhost:8000",
			ClientID:      "test-client",
			RedirectURL:   "http://localhost:8080/callback",
		}
		handler := NewAuthHandler(cfg)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/callback?state=invalid-state", nil)

		// Set different state in cookie
		c.Request.AddCookie(&http.Cookie{
			Name:  "oauth_state",
			Value: "valid-state",
		})

		handler.Callback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
