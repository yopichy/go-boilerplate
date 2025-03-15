package middleware

import (
	"github.com/gin-gonic/gin"
)

func Localization() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en" // Default language
		}
		c.Set("language", lang)
		c.Next()
	}
}
