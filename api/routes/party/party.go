package party

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Party struct {
	Name     string
	ServerID string
}

func PostParty(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	log.Info().
		Msg(string(body[:]))

	if err != nil {
		log.Err(err)
	}
	c.IndentedJSON(http.StatusOK, "Success!")
}
