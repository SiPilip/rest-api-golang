package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Hanya GetEvents yang pakai timeoutmiddleware
// @Summary      Get all events
// @Description  Fetch paginated list of events with optional search
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        page   query  int     false  "Page number"    default(1)
// @Param        limit  query  int     false  "Items per page" default(10)
// @Param        search query  string  false  "Search keyword"
// @Success      200    {object} helpers.PaginatedResponse
// @Failure      500    {object} helpers.Response
// @Router       /events [get]
func getEvents(context *gin.Context) {
	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(context.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	search := context.DefaultQuery("search", "")
	ctx := context.Request.Context()

	// 1. Cek cache
	cacheKey := fmt.Sprintf("events:page=%d:limit=%d:search=%s", page, limit, search)

	var cachedResponse helpers.PaginatedResponse
	if err := helpers.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		context.JSON(http.StatusOK, cachedResponse)
		return
	}

	// 2. Cache miss - Query Database
	events, total, err := models.GetAllEvents(page, limit, search, context.Request.Context())
	if err != nil {
		// Error timeout cause
		if context.Request.Context().Err() != nil {
        return
    }

		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not fetch events.")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// helpers.SuccessPaginatedResponse(context, http.StatusOK, "Events fetched succesfully.", events, helpers.Meta{
	// 	Page: page,
	// 	Limit: limit,
	// 	Total: total,
	// 	TotalPages: totalPages,
	// })

	response := helpers.PaginatedResponse{
		Status: "success",
		Message: "Events fetched succesfully.",
		Data: events,
		Meta: helpers.Meta{
			Page: page,
			Limit: limit,
			Total: total,
			TotalPages: totalPages,
		},
	}

	// 3. Cache miss, langsung simopan ke cache (TTL 30 detik)
	helpers.SetCache(ctx, cacheKey, response, 30*time.Second)
	context.JSON(http.StatusOK, response)
}

// @Summary      Get event by ID
// @Description  Fetch a single event by its ID
// @Tags         Events
// @Produce      json
// @Param        id  path  int  true  "Event ID"
// @Success      200 {object} helpers.Response
// @Failure      404 {object} helpers.Response
// @Router       /events/{id} [get]
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

// @Summary      Create event
// @Description  Create a new event (supports multipart form with image upload)
// @Tags         Events
// @Accept       multipart/form-data
// @Produce      json
// @Param        name        formData string true  "Event name"
// @Param        description formData string true  "Event description"
// @Param        location    formData string true  "Event location"
// @Param        datetime    formData string true  "Event datetime (RFC3339)"
// @Param        image       formData file   false "Event image"
// @Success      201 {object} helpers.Response
// @Failure      400 {object} helpers.Response
// @Failure      500 {object} helpers.Response
// @Security     BearerAuth
// @Router       /events [post]
func createEvent(context *gin.Context) {
    var event models.Event
    err := context.ShouldBind(&event)
    if err != nil {
        helpers.ValidationErrorResponse(context, err)
        return
    }
    // Handle file upload
    file, err := context.FormFile("image")
    if err == nil {
        if validErr := helpers.ValidateImage(file); validErr != nil {
            helpers.ErrorResponse(context, http.StatusBadRequest, validErr.Error())
            return
        }
        filePath, saveErr := helpers.SaveUploadedFile(file, "uploads")
        if saveErr != nil {
            helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save image.")
            return
        }
        saveErr = context.SaveUploadedFile(file, filePath)
        if saveErr != nil {
            helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save image.")
            return
        }
        event.ImageURL = "/" + filePath
    }

    event.UserID = context.GetInt64("userId")
    err = event.Save()
    if err != nil {
        helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save event.")
        return
    }

		// Hapus cache redis setelah di save
		helpers.DeleteCacheByPattern(context.Request.Context(), "events:*")
		
    helpers.SuccessResponse(context, http.StatusCreated, "Event created successfully", event)
}

// @Summary      Update event
// @Description  Update an existing event (owner only)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Event ID"
// @Success      200 {object} helpers.Response
// @Failure      400 {object} helpers.Response
// @Failure      401 {object} helpers.Response
// @Security     BearerAuth
// @Router       /events/{id} [put]
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
	role := context.GetString("role")
	if event.UserID != userId && role != "admin"{
		helpers.ErrorResponse(context, http.StatusForbidden, "You don't have permission to modify this event.")
		return
	}
	
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		helpers.ValidationErrorResponse(context,err)
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not update event!")
		return
	}

	// Hapus cache redis setelah di save
	helpers.DeleteCacheByPattern(context.Request.Context(), "events:*")

	helpers.SuccessResponse(context, http.StatusOK, "Event updated successfully", updatedEvent)
}

// @Summary      Delete event (soft delete)
// @Description  Soft delete an event (owner only)
// @Tags         Events
// @Produce      json
// @Param        id  path  int  true  "Event ID"
// @Success      200 {object} helpers.Response
// @Failure      401 {object} helpers.Response
// @Security     BearerAuth
// @Router       /events/{id} [delete]
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
	
	role := context.GetString("role")
	if event.UserID != userId && role != "admin" {
		helpers.ErrorResponse(context, http.StatusForbidden, "You don't have permission to delete this event.")
		return
	}

	err = event.Delete()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not delete event.")
		return
	}

	// Hapus cache redis setelah di save
	helpers.DeleteCacheByPattern(context.Request.Context(), "events:*")

	helpers.SuccessResponse(context, http.StatusOK, "Event deleted successfully.", nil)
}