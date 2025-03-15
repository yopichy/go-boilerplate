package middleware

import (
	"fmt"
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

		fmt.Printf("Debug - Found token: %s\n", token) // Add debug logging

		// Add token to request header if not present
		if c.GetHeader("Authorization") == "" {
			c.Request.Header.Set("Authorization", "Bearer "+token)
		}

		c.Next()
	}
}

func extractTokenFromMultipleSources(c *gin.Context) string {
	// Check fragment/hash for access_token (Swagger UI implicit flow)
	referer := c.GetHeader("Referer")
	if strings.Contains(referer, "access_token=") {
		parts := strings.Split(referer, "access_token=")
		if len(parts) > 1 {
			token := strings.Split(parts[1], "&")[0]
			return token
		}
	}

	// Check Authorization header
	if auth := c.GetHeader("Authorization"); auth != "" {
		parts := strings.Split(auth, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Check query parameter
	if token := c.Query("access_token"); token != "" {
		return token
	}

	return ""
}
