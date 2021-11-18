package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type api struct {
	tx *gorm.DB
}

func RegisterApi(router *gin.RouterGroup, tx *gorm.DB) {
	api := api{tx}
	router.POST("/register", api.registerUrl)
}

func (api *api) registerUrl(g *gin.Context) {

}
