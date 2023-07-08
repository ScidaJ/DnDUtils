package routes

import (
	"dndutils/api/controllers"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func UserRoute(router *gin.Engine, log *zerolog.Logger) {
	router.POST("/user", logger.SetLogger(), controllers.CreateUser())
	router.GET("/user", logger.SetLogger(), controllers.GetAllUsers())
}
