package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sayedulkrm/golang-jwt-project/controllers"
	// "github.com/sayedulkrm/golang-jwt-project/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	// Check if user is authenticated or not authenticated
	// incomingRoutes.Use(middlewares.AuthMiddlewares())

	incomingRoutes.GET("/users")                                                //get all users
	incomingRoutes.GET("/users/:user_id")                                       //get single user
	incomingRoutes.GET("/admin/users/:user_id", controllers.AdminGetAllUsers()) //get single user by id === Admin

}
