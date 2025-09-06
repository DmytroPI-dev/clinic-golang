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
	if err := db.AutoMigrate(&models.Program{}); err != nil {
		log.Fatalf("migration failed: %s", err)
	}
	log.Println("Migration successful")

	// Creating Gin router
	router := gin.Default()
	// Grouping API routes under /api/v1
	v1 := router.Group("/api/v1")
	{
		// Programs list endpoint
		v1.GET("/programs", handler.ListPrograms(db))
		// Program detail endpoint
		v1.GET("/programs/:id", handler.GetProgram(db))
		// Program creation endpoint
		v1.POST("/programs", handler.CreateProgram(db))
		//Program update endpoint
		v1.PUT("/programs/:id", handler.UpdateProgram(db))
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
