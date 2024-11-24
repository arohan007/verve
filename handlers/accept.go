package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"projects/verve/utils"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Global variables to track unique requests
var (
	uniqueIDs       = make(map[int]struct{})
	uniqueCount     int
	uniqueCountLock sync.Mutex
)

// Data structure for the POST request payload
type EndpointPayload struct {
	UniqueCount int    `json:"unique_count"`
	Message     string `json:"message"`
}

// HandleAccept processes the /api/verve/accept endpoint
func HandleAccept(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		idParam := c.Query("id")
		endpoint := c.Query("endpoint")

		// Validate id parameter
		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Println("Failed request: invalid id parameter")
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id parameter"})
			return
		}

		// Use Redis to ensure deduplication
		ctx := context.Background()
		redisKey := "unique_ids"
		added, err := utils.RedisClient.SAdd(ctx, redisKey, id).Result()
		if err != nil {
			logger.Printf("Error accessing Redis: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "internal server error"})
			return
		}

		if added == 0 {
			// ID already exists
			logger.Printf("Duplicate ID detected: %d\n", id)
		} else {
			logger.Printf("Unique ID added: %d\n", id)
		}

		// Optional: Send a POST request to the provided endpoint
		if endpoint != "" {
			go fireEndpointRequest(endpoint, logger)
		}

		// Respond with "ok"
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// StartUniqueCountTracker logs and resets the unique request count every minute
func StartUniqueCountTracker(logger *log.Logger) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		redisKey := "unique_ids"

		// Get unique count
		uniqueCount, err := utils.RedisClient.SCard(ctx, redisKey).Result()
		if err != nil {
			logger.Printf("Error accessing Redis: %v\n", err)
			continue
		}

		// Log the unique count
		logger.Printf("Unique requests in the last minute: %d\n", uniqueCount)

		// Reset the set
		_, err = utils.RedisClient.Del(ctx, redisKey).Result()
		if err != nil {
			logger.Printf("Error resetting Redis key: %v\n", err)
		}
	}
}


// fireEndpointRequest sends a POST request to the provided endpoint with unique count in the body
func fireEndpointRequest(endpoint string, logger *log.Logger) {
	uniqueCountLock.Lock()
	currentCount := uniqueCount
	uniqueCountLock.Unlock()

	// Prepare the payload
	payload := EndpointPayload{
		UniqueCount: currentCount,
		Message:     "Unique request count for the current minute",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	// Send the POST request
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Printf("Failed to send POST request to endpoint: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Log the response status
	logger.Printf("POST request sent to %s, response status: %d\n", endpoint, resp.StatusCode)
}
