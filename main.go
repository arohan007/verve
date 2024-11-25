package main

import (
	"os"
	"projects/verve/handlers"
	"projects/verve/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger := utils.InitLogger("requests.log")

	// Read environment variables for Redis and Kafka addresses
	redisAddr := os.Getenv("REDIS_ADDRESS")
	kafkaBrokerAddr := os.Getenv("KAFKA_BROKER_ADDRESS")

	// Initialize Redis
	utils.InitializeRedis(redisAddr) // Replace with your Redis address

	// Initialize Kafka writer
	kafkaWriter := utils.InitializeKafkaWriter(kafkaBrokerAddr, "unique-ids-count")
	defer kafkaWriter.Close()

	// Start unique request tracking in the background
	go handlers.StartUniqueCountTracker(logger, kafkaWriter)

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
