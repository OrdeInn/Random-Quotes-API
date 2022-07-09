package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"orderinn.com/random-quotes/configs"
	"orderinn.com/random-quotes/models"
	"orderinn.com/random-quotes/responses"
	"time"
)

var quoteCollection *mongo.Collection = configs.GetCollection(configs.DB, "quotes")

func getLastDocumentId() int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var lastQuote models.Quotation
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"_id", -1}})

	err := quoteCollection.FindOne(ctx, bson.D{}, findOptions).Decode(&lastQuote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0
		}
		return -1
	}

	return lastQuote.Id
}

func CreateQuote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var quote models.Quotation
		defer cancel()

		err := c.BindJSON(&quote)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				responses.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})

			return
		}

		lastQuoteId := getLastDocumentId()

		if lastQuoteId > -1 {
			lastQuoteId += 1
		}

		newQuote := models.Quotation{
			Id:        lastQuoteId,
			Author:    quote.Author,
			Quotation: quote.Quotation,
		}

		result, err := quoteCollection.InsertOne(ctx, newQuote)

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})

			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": result}})
	}
}

func GetRandomQuote() gin.HandlerFunc {
	return func(c *gin.Context) {

		lastQuoteId := getLastDocumentId()

		if lastQuoteId < 2 {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError,
					Message: "There is no enough quotation", Data: map[string]interface{}{"data": "error"}})

			return
		}

		randId := 0

		for randId == 0 {
			randId = rand.Intn(lastQuoteId + 1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var quote models.Quotation
		defer cancel()

		err := quoteCollection.FindOne(ctx, bson.M{"id": randId}).Decode(&quote)

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})

			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": quote}})
	}
}
