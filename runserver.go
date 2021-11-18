package main

import (
	"github.com/ddddanil/i.go/api"
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

func main() {
	db, err := initDb()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	gApi := router.Group("/api")
	api.RegisterApi(gApi, db)
	log.Fatalln(router.Run(":8080"))
}
