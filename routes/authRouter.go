package routes

import (
	"github.com/TheMikeKaisen/go_JWT/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(authRouter *gin.Engine) {

	authRouter.POST("/users/signup", controllers.Signup())
	authRouter.POST("/users/login", controllers.Login())

}
