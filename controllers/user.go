package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"orderinn.com/random-quotes/configs"
	"orderinn.com/random-quotes/models"
	"orderinn.com/random-quotes/responses"
	"time"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {

		var userInput RegisterInput
		var newUser models.User

		err := c.ShouldBindJSON(&userInput)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser.Username = userInput.Username
		newUser.Password, err = bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		result, err := userCollection.InsertOne(ctx, newUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"data": result,
			},
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
