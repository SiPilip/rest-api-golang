package routes

import (
	"REST-API/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	// Cek database
	dbStatus := "up"
	if err := db.DB.Ping(); err != nil {
		dbStatus = "down"
	}

	// Cek redis
	redisStatus := "up"
	if db.RedisClient == nil {
		redisStatus = "down"
	} else if err := db.RedisClient.Ping(c.Request.Context()).Err(); err != nil {
		redisStatus = "down"
	}

	// Overall status
	status := "healthy"
	httpCode := http.StatusOK
	if dbStatus == "down" {
		status = "unhealthy"
		httpCode = http.StatusServiceUnavailable
	}

	c.JSON(httpCode, gin.H{
		"status": status,
		"service": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	})
}