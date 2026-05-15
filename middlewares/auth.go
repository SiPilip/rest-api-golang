package middlewares

import (
	"REST-API/helpers"
	"REST-API/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	tokenString := context.Request.Header.Get("Authorization")
	if tokenString == "" {
		helpers.ErrorAuthResponse(context, http.StatusUnauthorized, "Server error. Authentication is required or has failed.")
		context.Abort()
		return
	}

	userId, role, err := utils.VerifyToken(tokenString)
	if err != nil {
		helpers.ErrorAuthResponse(context, http.StatusUnauthorized, "Server error. Authentication is required or has failed.")
		context.Abort()
		return
	}

	context.Set("userId", userId)
	context.Set("role", role)
	context.Next()
}