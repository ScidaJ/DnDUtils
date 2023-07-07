package main

import (
	"dndutils/api/routes/party"
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := gin.Default()

	r.POST("/party", logger.SetLogger(), party.PostParty)

	r.Run("localhost:8080")
}
