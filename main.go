package main

import (
	"projects/verve/handlers"
	"projects/verve/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger := utils.InitLogger("requests.log")

	// Start unique request tracking in the background
	go handlers.StartUniqueCountTracker(logger)

	// Create a Gin router
	r := gin.Default()

	// Define routes
	r.GET("/api/verve/accept", handlers.HandleAccept(logger))

	// Start server
	logger.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
