package routes

import (
	"dndutils/api/controllers"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func PartyRoute(router *gin.Engine, log *zerolog.Logger) {
	router.POST("/party", logger.SetLogger(), controllers.CreateParty())
}
