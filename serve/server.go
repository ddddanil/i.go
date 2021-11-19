package serve

import (
	"context"
	"errors"
	"github.com/ddddanil/i.go/api"
	"github.com/ddddanil/i.go/html"
	"github.com/ddddanil/i.go/shortener"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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
	router.SetHTMLTemplate(html.TemplateHtml())
	apiGroup := router.Group("/api")
	api.RegisterApi(apiGroup)
	router.GET("/", html.IndexHandler)
	router.GET("/index", html.IndexHandler)
	router.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		c.Request.Form = url.Values{"short": []string{short}}
		c.Request.URL.Path = "/api/redirect"
		router.HandleContext(c)
	})
	return router
}

func SystemQuit() <-chan os.Signal {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func startServer(h http.Handler) (srv *http.Server, e *errgroup.Group, ctx context.Context) {
	srv = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	e, ctx = errgroup.WithContext(context.Background())

	// Initializing the serve in a goroutine so that
	// it won't block the graceful shutdown handling below
	e.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		} else if err != nil {
			log.Fatal(err)
		}
		return err
	})
	return srv, e, ctx
}

func WaitAsync(e *errgroup.Group) <-chan error {
	result := make(chan error)
	go func() {
		err := e.Wait()
		result <- err
	}()
	return result
}

func ShutdownServer(srv *http.Server) {
	log.Println("Shutting down server...")

	// The context is used to inform the serve it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

func Run() {
	db, err := initDb()
	if err != nil {
		panic(err)
	}
	router := configRouter(db)
	srv, e, ctx := startServer(router)
	shortener.DeleteExpired(ctx, db)

	serverError := WaitAsync(e)
	systemQuit := SystemQuit()

	for {
		select {
		case err = <-serverError:
			if err != nil {
				log.Fatal(err)
			}
			return
		case <-systemQuit:
			ShutdownServer(srv)
			return
		}
	}
}
