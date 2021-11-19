package html

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexHTML struct {
	RegisterUrl string
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", IndexHTML{"/api/register"})
}
