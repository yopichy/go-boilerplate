package middleware

import (
	"webapi/internal/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors after request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error(err)
			}
		}
	}
}
