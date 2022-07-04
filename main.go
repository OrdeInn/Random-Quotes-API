package main

import(
	"orderinn.com/random-quotes/configs"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	configs.ConnectDB()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data" : "Hello World!",
		})
	})

	router.Run("localhost:8080")
}