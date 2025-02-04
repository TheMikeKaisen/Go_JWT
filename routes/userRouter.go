package routes

import (
	"github.com/TheMikeKaisen/go_JWT/controllers"
	"github.com/TheMikeKaisen/go_JWT/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(userRouter *gin.Engine) {
	userRouter.Use(middleware.Authenticate())
	userRouter.GET("/user/:userId", controllers.GetUser())
}
