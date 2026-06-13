package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"ancient-text-backend/internal/cache"
	"ancient-text-backend/internal/config"
	"ancient-text-backend/internal/database"
	"ancient-text-backend/internal/router"
)

func main() {
	config.Load()

	gin.SetMode(config.AppConfig.Server.Mode)

	database.Connect()

	cache.Connect()

	r := router.Setup()

	port := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on port %s", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
