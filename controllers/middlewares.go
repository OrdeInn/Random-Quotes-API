package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	token "orderinn.com/random-quotes/utils"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := token.TokenValid(c)

		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
