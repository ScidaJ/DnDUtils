package controllers

import (
	"context"
	"dndutils/api/configs"
	"dndutils/api/models"
	"dndutils/api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var partyCollection *mongo.Collection = configs.GetCollection(configs.DB, "parties")
var validate = validator.New()

func CreateParty() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var party models.Party

		//validate the request body
		if err := c.BindJSON(&party); err != nil {
			c.JSON(http.StatusBadRequest, responses.PartyResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use valiator to validate request fields
		if validationErr := validate.Struct(&party); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PartyResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		}

		newParty := models.Party{
			Id:       primitive.NewObjectID(),
			Name:     party.Name,
			ServerId: party.ServerId,
			Owner:    party.Owner,
		}

		result, err := partyCollection.InsertOne(ctx, newParty)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PartyResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		c.JSON(http.StatusCreated, responses.PartyResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
