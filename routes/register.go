package routes

import (
	"REST-API/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not register user for event."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registered successfully!"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not cancel registration."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Registration cancelled successfully!"})
}