package main

import (
	"dndutils/api/routes/party"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/party", party.PostParty)

	r.Run("localhost:8080")
}
