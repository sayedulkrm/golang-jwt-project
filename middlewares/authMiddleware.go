package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayedulkrm/golang-jwt-project/helpers"
)

func AuthMiddlewares() gin.HandlerFunc {

	return func(c *gin.Context) {

		accessTokenCookie, err := c.Request.Cookie("access_token")
		if err != nil {
			// Handle the error if the cookie is not found
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
			c.Abort() // Abort the request processing
			return
		}

		// Get the value of the access token from the cookie
		accessToken := accessTokenCookie.Value

		// Check if the access token is empty
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is empty"})
			c.Abort() // Abort the request processing
			return
		}

		claims, err := helpers.ValidateToken(accessToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort() // Abort the request processing
			return
		}

		fmt.Println("Claims:", claims)

		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("user_id", claims.UserId)
		c.Set("user_type", claims.User_type)

		c.Next()

	}

}
