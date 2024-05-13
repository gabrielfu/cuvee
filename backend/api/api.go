package api

import (
	"context"
	"cuvee/db"
	"cuvee/domain/images"
	"cuvee/domain/wines"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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

func initRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	return r
}

func Run() {
	r := initRouter()

	connector := db.NewMongoConnector(
		os.Getenv("MONGO_URI"),
		os.Getenv("MONGO_DATABASE"),
		os.Getenv("MONGO_COLLECTION"),
	)
	collection, err := connector.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := collection.Database().Client().Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// register wine service
	repo := wines.NewWineRepository(collection)
	validate := wines.NewWineJSONValidator()
	service := wines.NewWineService(repo, validate)
	wines.RegisterRoutes(r, service)

	// register image service
	imageService := images.NewImageService(os.Getenv("GOOGLE_SEARCH_CX"), os.Getenv("GOOGLE_SEARCH_API_KEY"))
	images.RegisterRoutes(r, imageService)

	initServer(r)
}
