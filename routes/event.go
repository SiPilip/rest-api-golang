package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not fetch events.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Events fetched succesfully!", events)
}

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not parse event id.")
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		helpers.ErrorResponse(context, http.StatusNotFound, "Event not found.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Event fetched successfully.", event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "CoCould not parse request data.")
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save event.")
		return
	}

	helpers.SuccessResponse(context, http.StatusCreated, "Event created successfully.", event)
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse event id.")
		return
	}

	
	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusNotFound, "Could not find event by id.")
		return
	}
	if event.UserID != userId {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Client error. User can't update other user's event.")
		return
	}
	
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse request data!")
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not update event!")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Event updated successfully", updatedEvent)
}

func deleteEvent(context *gin.Context) {
	// CHECK IF HAS PARAM
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse event id.")
		return
	}

	// GET EVENT BY ID
	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusNotFound, "Could not find event by id.")
		return
	}
	
	if event.UserID != userId {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Client error. User can't update other user's event.")
		return
	}

	err = event.Delete()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not delete event.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Event deleted successfully.", nil)
}