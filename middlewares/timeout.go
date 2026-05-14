package middlewares

import (
	"REST-API/helpers"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Buat context dengan timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Ganti request context dengan yang ada timeout
		c.Request = c.Request.WithContext(ctx)

		// Lanjut proses request
		c.Next()

		// Cek apakah timeout terjadi
		if ctx.Err() == context.DeadlineExceeded {
			helpers.ErrorResponse(c, http.StatusGatewayTimeout, "Request timeout.")
			c.Abort()
		}
	}
}

// context.WithTimeout — membuat "stopwatch" yang otomatis expired setelah durasi tertentu. Semua operasi yang menggunakan context ini (termasuk database query) akan tahu bahwa waktunya sudah habis.