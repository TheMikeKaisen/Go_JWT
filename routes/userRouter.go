package routes

import (
	"github.com/TheMikeKaisen/go_JWT/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(userRouter *gin.Engine) {
	// userRouter.Use(middleware.AuthMiddleware())
	// userRouter.GET("/users", controllers.GetUsers())
	userRouter.GET("/user/:userId", controllers.GetUser())
}
