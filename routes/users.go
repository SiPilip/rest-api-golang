package routes

import (
	"REST-API/helpers"
	"REST-API/models"
	"REST-API/utils"
	"REST-API/workers"
	"net/http"
	"time"

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
	if err != nil {
		helpers.ValidationErrorResponse(context,err)
		return
	}

	err = user.Save()

	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save user!")
		return
	}

	workers.Enqueue(workers.Job{
    Name: "send_welcome_email",
    Execute: func() error {
        return helpers.SendEmail(helpers.EmailData{
            To:      user.Email,
            Subject: "Welcome to Event API!",
            Body:    helpers.WelcomeEmailBody(user.Email),
        })
    },
})

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

	// Generate access token (15 minutes)
	accessToken, err := utils.GenerateAccessToken(user.Email, user.ID, user.Role)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, "Client error. Authentication is required or has failed.")
		return
	}

	// Generate refresh token (7 hari)
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not generate token.")
		return
	}

	// Simpan refresh token di database
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = models.SaveRefreshToken(user.ID, refreshToken, expiresAt)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not save refresh token.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Login successful!", 
		gin.H{
			"access_token": accessToken,
			"refresh_token": refreshToken,
			"token_type": "Bearer",
			"expires_in": 15 * 60, // 15 minutes in seconds
			"user": gin.H{
				"id": user.ID,
				"email": user.Email,
				"role": user.Role,
			},
		},
	)
}	

func refreshAccessToken(context *gin.Context) {
	var body struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
  }

	if err := context.ShouldBindJSON(&body); err != nil {
		helpers.ValidationErrorResponse(context, err)
		return
	}

	// Cari refresh token  di database
	rt, err := models.GetRefreshToken(body.RefreshToken)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Invalid refresh token.")
		return
	}

	// Cek apakah expired
	if time.Now().After(rt.ExpiresAt) {
		models.DeleteRefreshToken(body.RefreshToken)
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Refresh token expired. Please login again.")
		return
	}

	// Ambil data user untuk generate token baru
	user, err := models.GetUserByID(rt.UserID)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not find user.")
		return
	}

	// Generate access token baru
	accessToken, err := utils.GenerateAccessToken(user.Email, user.ID, user.Role)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not generate token.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Token refreshed", 
		gin.H {
			"access_token": accessToken,
		},
	)
}

func logout(context *gin.Context){
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		helpers.ValidationErrorResponse(context, err)
		return
	}

	err := models.DeleteRefreshToken(body.RefreshToken)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusInternalServerError, "Could not logout.")
		return
	}

	helpers.SuccessResponse(context, http.StatusOK, "Logged out successfully.", nil)
}