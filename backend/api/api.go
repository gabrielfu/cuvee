package api

import (
	"context"
	"cuvee/db"
	"cuvee/domain/images"
	"cuvee/domain/regions"
	vintagecharts "cuvee/domain/vintage_charts"
	"cuvee/domain/wines"
	"cuvee/external/llm"
	"cuvee/external/search"
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
	)
	db, err := connector.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := db.Client().Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// register wine service
	wineCollection := db.Collection(os.Getenv("MONGO_WINE_COLLECTION"))
	wineRepo := wines.NewWineRepository(wineCollection)
	validate := wines.NewWineJSONValidator()
	wineService := wines.NewWineService(wineRepo, validate)
	wines.RegisterRoutes(r, wineService)

	// register image service
	searchEngine, err := search.NewGoogleSearchEngine(os.Getenv("GOOGLE_SEARCH_CX"), os.Getenv("GOOGLE_SEARCH_API_KEY"))
	if err != nil {
		log.Fatalf("Failed to create Google search engine: %v", err)
	}
	imageService := images.NewImageService(searchEngine)
	images.RegisterRoutes(r, imageService)

	// register region service
	regionCollection := db.Collection(os.Getenv("MONGO_REGION_COLLECTION"))
	regionRepo := regions.NewRegionRepository(context.Background(), regionCollection)
	regionService := regions.NewRegionService(regionRepo)
	regions.RegisterRoutes(r, regionService)

	// register vintage charts service
	rpProvider, err := vintagecharts.NewRPProvider(
		"data/vintage_charts/robert_parker_chart.json",
		"data/vintage_charts/robert_parker_maturity.json",
	)
	if err != nil {
		log.Fatalf("Failed to create RP provider: %v", err)
	}
	providers := []vintagecharts.Provider{rpProvider}
	openaiLLM := llm.NewOpenAILLM(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
	googleSearch, err := search.NewGoogleSearchEngine(os.Getenv("GOOGLE_SEARCH_CX"), os.Getenv("GOOGLE_SEARCH_API_KEY"))
	if err != nil {
		log.Fatalf("Failed to create Google search engine: %v", err)
	}
	vintageChartService := vintagecharts.NewVintageChartService(providers, openaiLLM, googleSearch)
	vintagecharts.RegisterRoutes(r, vintageChartService)

	initServer(r)
}
