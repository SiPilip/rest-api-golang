package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"REST-API/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User credentials"
// @Success      201  {object} helpers.Response
// @Failure      400  {object} helpers.Response
// @Router       /signup [post]
func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	// if err != nil {
	// 	helpers.ErrorResponse(context, http.StatusBadRequest, "Could not parse request data!")
	// 	return
	// }
	if err != nil {
		helpers.ValidationErrorResponse(context,err)
		return
	}

	err = user.Save()

	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save user!")
		return
	}

	helpers.SuccessResponse(context, http.StatusCreated, "User created successfully", nil)
}

// @Summary      Login
// @Description  Authenticate user and get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User credentials"
// @Success      200  {object} helpers.Response
// @Failure      401  {object} helpers.Response
// @Router       /login [post]
func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		helpers.ValidationErrorResponse(context,err)
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