package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Request proccess
		c.Next()

		// Log setelah request selesai	
		duration := time.Since(start)

		slog.Info("Request",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"duration", duration.String(),
				"ip", c.ClientIP())
	}
}