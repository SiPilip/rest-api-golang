package helpers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data any) {
	// context.JSON(http.StatusOK, events)
	ctx.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	// context.JSON(http.StatusInternalServerError, gin.H{
	// 	"message": "Could not fetch events.",
	// })
	ctx.JSON(statusCode,Response{
		Status: "error",
		Message: message,
	})
}

func ErrorAuthResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode,Response{
		Status: "error",
		Message: message,
	})
}

func ValidaitonErrorResponse(ctx *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		messages := make([]string, 0, len(validationErrors))
		for _, e := range validationErrors {
			switch e.Tag(){
				case "required":
					messages = append(messages, e.Field()+" is required.")
				case "email":
					messages = append(messages, e.Field()+" must be a valid email.")
				case "min":
					messages = append(messages, e.Field()+" must be a at least " + e.Param()+" characters.")
				case "max":
					messages = append(messages, e.Field()+" must be a at most " + e.Param()+" characters.")
				default:
					messages = append(messages, e.Field()+" is invalid.")
			}
		}
		ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Validation failed.",
			Data: messages,
		})
		return
	}
	
	// incase it was not a validation error.
	ErrorResponse(ctx, http.StatusBadRequest, "Invalid request data")
}