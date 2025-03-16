package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := "en" // default language
		acceptLang := c.GetHeader("Accept-Language")

		if acceptLang != "" {
			// Get the first language from Accept-Language header
			langs := strings.Split(acceptLang, ",")
			if len(langs) > 0 {
				// Extract language code (e.g., "en" from "en-US")
				lang = strings.Split(langs[0], "-")[0]
			}
		}

		c.Set("language", lang)
		c.Next()
	}
}
