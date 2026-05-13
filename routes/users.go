package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"REST-API/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	// if err != nil {
	// 	helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse request data!")
	// 	return
	// }
	if err != nil {
		helpers.ValidaitonErrorResponse(context,err)
		return
	}

	err = user.Save()

	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save user!")
		return
	}

	helpers.SuccessResponse(context, http.StatusCreated, "User created successfully", nil)
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		helpers.ValidaitonErrorResponse(context,err)
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Client error. Authentication is required or has failed.")
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Client error. Authentication is required or has failed.")
		return
	}

	data := gin.H{
		"token": token,
	}
	helpers.SuccessResponse(context, http.StatusOK, "Login successful!", data)
}