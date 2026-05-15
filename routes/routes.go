package routes

import (
	"REST-API/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Register user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 201 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Router /signup [post]
func RegisterRoutes(server *gin.Engine) {
	server.Use(middlewares.RateLimiter())
	server.Use(middlewares.CORSMiddleware())
	server.Use(middlewares.TimeoutMiddleware(5 * time.Second))
	server.Use(middlewares.RequestLogger())

	// API v1
	v1 := server.Group("/api/v1")
	{
		v1.GET("/events", getEvents)
		v1.GET("/events/:id", getEvent)

		authenticated := v1.Group("/")
		authenticated.Use(middlewares.Authenticate)
		authenticated.POST("/events", createEvent)
		authenticated.PUT("/events/:id", updateEvent)
		authenticated.DELETE("/events/:id", deleteEvent)
		authenticated.POST("/events/:id/register", registerForEvent)
		authenticated.DELETE("/events/:id/register", cancelRegistration)

		v1.POST("/signup", signup)
		v1.POST("/login", login)
		v1.POST("/refresh", refreshAccessToken)
		v1.POST("/logout", logout)
	}
	// Admin only routes
	admin := v1.Group("/admin")
	admin.Use(middlewares.Authenticate)
	admin.Use(middlewares.RequireRole("admin"))
	{
			// Contoh: endpoint untuk lihat semua users (nanti)
			// admin.GET("/users", getAllUsers)
	}

	// API v2 — contoh di masa depan
	// v2 := server.Group("/api/v2")
	// {
	// 	v2.GET("/events", getEventsV2)  // handler baru dengan format response berbeda
	// 	// ... route v2 lainnya
	// }
}