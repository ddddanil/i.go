package api

import (
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type api struct {
	tx *gorm.DB
}

func RegisterApi(router *gin.RouterGroup, tx *gorm.DB) {
	api := api{tx}
	router.POST("/register", api.registerUrl)
}

type urlView struct {
	url      string `binding:"required"`
	expireIn int
}

func (api *api) registerUrl(c *gin.Context) {
	var url urlView
	if err := c.ShouldBind(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var short shortener.ShortUrl
	if url.expireIn != 0 {
		short = shortener.NewShortUrl(url.url,
			shortener.WithExpiration(time.Duration(url.expireIn)*time.Minute))
	} else {
		short = shortener.NewShortUrl(url.url)
	}
	api.tx.Create(&short)
	c.JSON(http.StatusOK, gin.H{
		"short": short.Shortened,
	})
}
