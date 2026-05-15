package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cek client sudah kirim requestID atau belum
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Simpan di context (dipakai di handler)
		c.Set("requestId", requestID)

		// Tambahkan di response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}