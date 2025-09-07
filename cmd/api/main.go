package main

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/config"
	"github.com/DmytroPI-dev/clinic-golang/internal/database"
	"github.com/DmytroPI-dev/clinic-golang/internal/handler"
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	//Load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load environment variables: %s", err)
	}

	//Connect to DB
	db, err := database.DB_Connect(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	log.Println("Successfully connected to database")
	// Migrating data
	log.Println("Starting DB migration....")
	if err := db.AutoMigrate(&models.Program{}, &models.Price{}, &models.News{}, &models.User{}); err != nil {
		log.Fatalf("migration for models.Program failed: %s", err)
	}
	log.Println("Migration successful")

	// Creating Gin router
	router := gin.Default()
	// Grouping API routes under /api/v1
	v1 := router.Group("/api/v1")
	{
		{
			// Programs list endpoints
			programRoutes := v1.Group("/programs")
			{
				programRoutes.GET("/", handler.ListPrograms(db))
				programRoutes.GET("/:id", handler.GetProgram(db))
				programRoutes.POST("/", handler.CreateProgram(db))
				programRoutes.PUT("/:id", handler.UpdateProgram(db))
				programRoutes.DELETE("/:id", handler.DeleteProgram(db))
			}
			// Prices list endpoints
			priceRoutes := v1.Group("/prices")
			{
				priceRoutes.GET("/", handler.ListPrices(db))
				priceRoutes.GET("/:id", handler.GetPrice(db))
				priceRoutes.POST("/", handler.CreatePrice(db))
				priceRoutes.PUT("/:id", handler.UpdatePrice(db))
				priceRoutes.DELETE("/:id", handler.DeletePrice(db))
			}
			// News list endpoints
			newsRoutes := v1.Group("/news")
			{
				newsRoutes.GET("/", handler.ListNews(db))
				newsRoutes.GET("/:id", handler.GetNews(db))
				newsRoutes.POST("/", handler.CreateNews(db))
				newsRoutes.PUT("/:id", handler.UpdateNews(db))
				newsRoutes.DELETE("/:id", handler.DeleteNews(db))
			}
		}
	}

	//Testing
	router.GET("/ping", func(ctx *gin.Context) {
		//c.Json sends response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Start server
	serverAddress := "localhost:" + cfg.ServerPort
	log.Printf("Starting server on %s", serverAddress)
	router.Run(serverAddress)

}
