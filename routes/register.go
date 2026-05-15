package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse event id.")
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusNotFound, "Event not found.")
		return
	}

	err = event.Register(userId)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not register user for event.")
		return
	}

	go func () {
		// Ambil data user & event untuk email
		user, err := models.GetUserByID(userId)
		if err != nil {
			slog.Error("Failed to get user for email", "error", err)
			return
		}

		helpers.SendEmail(helpers.EmailData{
			To: user.Email,
			Subject: "Registration Confirmed: " + event.Name,
			Body: helpers.EventRegistrationEmailBody(user.Email, event.Name),
		})
	}()

	helpers.SuccessResponse(context, http.StatusCreated, "Registered successfully!", nil)
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse event id.")
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusNotFound, "Could not cancel registration.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Registration cancelled successfully!", nil)
}