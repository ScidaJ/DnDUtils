package main

import (
	"dndutils/api/configs"
	"dndutils/api/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	configs.ConnectDB()

	r := gin.Default()

	routes.PartyRoute(r, &log.Logger)
	routes.UserRoute(r, &log.Logger)

	r.Run("localhost:8080")
}
