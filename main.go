package main

import (
	"fmt"
	"os"

	"github.com/TheMikeKaisen/go_JWT/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)



func main() {

	// load env
	loadEnvErr := godotenv.Load(".env")
	if loadEnvErr != nil {
		fmt.Println("Error while loading env")
	}

	// load gin
	router := gin.Default()


	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api/v1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message":"api/v2 router running successfully!"})
	})

	router.Run(":"+os.Getenv("PORT"))


}