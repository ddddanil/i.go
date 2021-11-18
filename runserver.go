package main

import (
	"context"
	"github.com/ddddanil/i.go/api"
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func initDb() (db *gorm.DB, err error) {
	dsn := "host=localhost user=igo password=igo dbname=igo port=5432"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	// Migrate the schema
	err = db.AutoMigrate(&shortener.ShortUrl{})
	if err != nil {
		return
	}
	return
}

func UseDbContext(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeoutContext, stop := context.WithTimeout(context.Background(), time.Second)
		c.Set("DB", db.WithContext(timeoutContext))
		c.Next()
		stop()
	}
}

func main() {
	db, err := initDb()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(UseDbContext(db))
	apiGroup := router.Group("/api")
	api.RegisterApi(apiGroup, db)
	log.Fatalln(router.Run(":8080"))
}
