package middleware

import (
	"net/http"
	"strings"
	"webapi/config"

	"github.com/gin-gonic/gin"
)

func OAuth2Authentication(cfg config.OAuth2Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromMultipleSources(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		// Add token to request header if not present
		if c.GetHeader("Authorization") == "" {
			c.Request.Header.Set("Authorization", "Bearer "+token)
		}

		c.Next()
	}
}

func extractTokenFromMultipleSources(c *gin.Context) string {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		return authHeader
	}

	// Try cookie
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie
	}

	// Try query parameter
	if token := c.Query("access_token"); token != "" {
		return token
	}

	return ""
}
