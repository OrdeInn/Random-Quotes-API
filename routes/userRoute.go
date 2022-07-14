package routes

import (
	"github.com/gin-gonic/gin"
	"orderinn.com/random-quotes/controllers"
)

func UserRoute(router *gin.Engine) {

	protectedGroup := router.Group("/api")
	publicGroup := router.Group("/")

	initProtectedRoutes(protectedGroup)
	initPublicRoutes(publicGroup)
}

func initPublicRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/login", controllers.Login())
	routerGroup.POST("/register", controllers.Register())
}

func initProtectedRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.Use(controllers.AuthorizationMiddleware())

	routerGroup.GET("/random", controllers.GetRandomQuote())
	routerGroup.POST("/", controllers.CreateQuote())
}
