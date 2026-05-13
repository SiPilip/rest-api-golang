package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Domain mana yang boleh diakses
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},

		// HTTP Method apa yang diizinkan
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// Header yang boleh dikirim client
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		
		// Header response yang boleh  dibaca client
		ExposeHeaders: []string{"Content-Length"},

		// Izinkan cookies/auth headers
		AllowCredentials: true,

		// Max age 12 jam
		MaxAge: 12 * time.Hour, 
	})
}