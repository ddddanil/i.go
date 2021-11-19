package api

import (
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func RegisterApi(router *gin.RouterGroup) {
	router.POST("/register", registerUrl)
	router.GET("/get", getUrl)
}

type urlView struct {
	url      string `binding:"required"`
	expireIn int
}

func registerUrl(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
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
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&short)
		return tx.Error
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"short": short.Shortened,
	})
	return
}

type getView struct {
	short string `form:"short" binding:"required"`
}

func getUrl(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var form getView
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	short, err := shortener.GetShortUrl(form.short, db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "This link does not exist or is expired"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, short.Redirect)
}
