package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sayedulkrm/golang-jwt-project/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("users/register", controllers.Register())
	incomingRoutes.POST("users/login", controllers.Login())

}
