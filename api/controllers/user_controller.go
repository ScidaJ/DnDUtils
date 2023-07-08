package controllers

import (
	"context"
	"dndutils/api/configs"
	"dndutils/api/models"
	"dndutils/api/responses"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var reg *regexp.Regexp = regexp.MustCompile("^[0-9]*$")

func GetUser() gin.HandlerFunc {
	return getUser
}

func GetAllUsers() gin.HandlerFunc {
	return getAllUsers
}

func CreateUser() gin.HandlerFunc {
	return createUser
}

func getUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")

	userIdClean := reg.ReplaceAllString(userId, "")

	if len(userIdClean) != 18 {
		c.JSON(http.StatusBadRequest, responses.Response{Data: map[string]interface{}{"data": "This error should never come up in normal operation. Contact me on Discord, jscida, for assistance."}})
		return
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"discord_id": userIdClean}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, responses.Response{Message: "success", Data: map[string]interface{}{"data": user}})
}

func getAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()

	var results []*models.User
	cur, err := userCollection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		results = append(results, &user)
	}

	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, responses.Response{Message: "success", Data: map[string]interface{}{"data": results}})
}

func createUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	//validate the request body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	//use valiator to validate request fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.Response{Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	exists, err := userExists(user.DiscordId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if !exists {
		newUser := models.User{
			Id:        primitive.NewObjectID(),
			DiscordId: user.DiscordId,
			Servers:   user.Servers,
			Owns:      user.Owns,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		c.JSON(http.StatusCreated, responses.Response{Message: "success", Data: map[string]interface{}{"data": result}})
	} else {
		//TODO: Change to PUT request on the object
		c.JSON(http.StatusConflict, responses.Response{Message: "error", Data: map[string]interface{}{"data": "userId already in use"}})
	}
}

func userExists(discordId string) (bool, error) {
	discordIdClean := reg.ReplaceAllString(discordId, "")

	findOptions := options.FindOne()
	filter := bson.D{primitive.E{Key: "discord_id", Value: discordIdClean}}

	var results bson.M
	err := userCollection.FindOne(context.TODO(), filter, findOptions).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
