package main

import (
	"github.com/gin-gonic/gin"
	"orderinn.com/random-quotes/configs"
	"orderinn.com/random-quotes/routes"
)

func main() {

	router := gin.Default()

	configs.ConnectDB()

	routes.UserRoute(router)

	router.Run("localhost:8080")
}
