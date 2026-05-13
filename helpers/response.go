package helpers

import "github.com/gin-gonic/gin"

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