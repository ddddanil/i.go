package main

import (
	"context"
	"errors"
	"github.com/ddddanil/i.go/api"
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initDb() (db *gorm.DB, err error) {
	dsn := "host=localhost user=igo password=igo dbname=igo port=5432"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
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
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
		c.Set("DB", db.WithContext(timeoutContext))
		//c.Set("DB", db)
		c.Next()
	}
}

func configRouter(db *gorm.DB) http.Handler {
	router := gin.Default()
	router.Use(UseDbContext(db))
	apiGroup := router.Group("/api")
	api.RegisterApi(apiGroup)
	router.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		c.Request.Form = url.Values{"short": []string{short}}
		c.Request.URL.Path = "/api/redirect"
		router.HandleContext(c)
	})
	return router
}

func main() {
	global, cancelGlobal := context.WithCancel(context.Background())
	defer cancelGlobal()
	db, err := initDb()
	if err != nil {
		panic(err)
	}
	router := configRouter(db)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		err := srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		} else if err != nil {
			log.Fatal(err)
		}
	}()
	shortener.DeleteExpired(global, db)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(global, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
