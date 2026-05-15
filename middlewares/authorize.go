package middlewares

import (
	"REST-API/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		helpers.ErrorResponse(c, http.StatusForbidden, "You dont have permission to access this resource.")
		c.Abort()
	}
}