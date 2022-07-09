package routes

import (
	"github.com/gin-gonic/gin"
	"orderinn.com/random-quotes/controllers"
)

func UserRoute(router *gin.Engine) {

	router.GET("/random", controllers.GetRandomQuote())

	router.POST("/", controllers.CreateQuote())
}
