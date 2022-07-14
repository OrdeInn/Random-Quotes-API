package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"orderinn.com/random-quotes/configs"
	"orderinn.com/random-quotes/models"
	"orderinn.com/random-quotes/responses"
	token "orderinn.com/random-quotes/utils"
	"time"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
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
		var input LoginInput

		//Validate input
		err := c.ShouldBindJSON(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Query the user from collection
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var user models.User

		err = userCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		//Validate password
		err = verifyPassword(&user, input.Password)

		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		userToken, err := token.GenerateToken(user.Id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"username": user.Username,
					"token":    userToken,
				},
			},
		})

	}
}

func verifyPassword(user *models.User, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))
}
