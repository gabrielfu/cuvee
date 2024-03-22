package api

import (
	"context"
	"cuvee/domain/wines"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func initServer(r *gin.Engine) {
	// graceful shutdown from Gin example
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func initRouter(service *wines.WineService) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/wines/:id", handleGetWine(service))
	r.GET("/wines", handleListWines(service))
	r.POST("/wines", handleCreateWine(service))
	r.PATCH("/wines/:id", handleUpdateWine(service))
	r.DELETE("/wines/:id", handleDeleteWine(service))

	return r
}

func Run() {
	service := &wines.WineService{}
	r := initRouter(service)
	initServer(r)
}