package party

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostParty(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Success!")
}
