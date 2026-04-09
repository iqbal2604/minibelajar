// main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mini/app/router"
	"mini/config"
	"mini/database/connection"
	_ "mini/database/ent/runtime"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc

	gin.SetMode(config.App.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	db := connection.DB()
	if db == nil {
		os.Exit(1)
	}

	// Auto migrate
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	rh := &router.Handlers{DB: db, R: r}
	rh.Routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.App.Port),
		Handler: r.Handler(),
	}

	go func() {
		log.Printf("Server running on :%s", config.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	connection.CloseDB()
	log.Println("Server exited")
}
