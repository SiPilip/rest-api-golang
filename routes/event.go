package routes

import (
	"REST-API/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events.",
		})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found.",
			"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request data!" + err.Error()})
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save event.",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created successfully!", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	
	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Could not find event by id.",
		})
		return
	}
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Client error. User can't update other user's event.",
		})
		return
	}
	
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request data!" + err.Error()})
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update event!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!", "event": updatedEvent})
}

func deleteEvent(context *gin.Context) {
	// CHECK IF HAS PARAM
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	// GET EVENT BY ID
	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Could not find event by id.",
		})
		return
	}
	
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Client error. User can't update other user's event.",
		})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete event!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}