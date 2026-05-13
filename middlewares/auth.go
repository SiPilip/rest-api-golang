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
		return
	}

	userId, err := utils.VerifyToken(tokenString)
	if err != nil {
		helpers.ErrorAuthResponse(context, http.StatusUnauthorized, "Server error. Authentication is required or has failed.")
		return
	}

	context.Set("userId", userId)
	context.Next()
}