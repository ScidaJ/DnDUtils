package party

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func postParty(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}
