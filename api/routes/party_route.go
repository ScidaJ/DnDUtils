package routes

import (
	"dndutils/api/controllers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func PartyRoute(router *gin.Engine, log *zerolog.Logger) {
	router.POST("/party", controllers.CreateParty())
}
