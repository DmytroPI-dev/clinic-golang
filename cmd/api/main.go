package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DmytroPI-dev/clinic-golang/internal/config"
	"github.com/DmytroPI-dev/clinic-golang/internal/database"
	"github.com/gin-gonic/gin"
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
	fmt.Println(db)


	// Creating Gin router
	router := gin.Default()
	//Testing
	router.GET("/ping", func(ctx *gin.Context) {
		//c.Json sends response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	serverAddress := "localhost:" + cfg.ServerPort
	log.Printf("Starting server on %s", serverAddress)
	router.Run(serverAddress)

}
