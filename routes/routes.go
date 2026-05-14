package routes

import (
	"REST-API/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.Use(middlewares.RateLimiter())
	server.Use(middlewares.CORSMiddleware())
	server.Use(middlewares.TimeoutMiddleware(5 * time.Second))
	server.Use(middlewares.RequestLogger())
	
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)
}