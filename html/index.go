package html

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type IndexHTML struct {
	RegisterUrl string
}

func IndexPage() *template.Template {
	return template.Must(template.ParseFiles("html/templates/index.html"))
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", IndexHTML{"/api/register"})
}
