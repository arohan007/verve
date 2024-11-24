package handlers

import (
	"log"
	"net/http"
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

		// Update unique request count
		uniqueCountLock.Lock()
		if _, exists := uniqueIDs[id]; !exists {
			uniqueIDs[id] = struct{}{}
			uniqueCount++
		}
		uniqueCountLock.Unlock()

		// Optional: Send a GET request to the provided endpoint
		if endpoint != "" {
			go fireEndpointRequest(endpoint, logger)
		}

		// Respond with "ok"
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// StartUniqueCountTracker logs and resets the unique request count every minute
func StartUniqueCountTracker(logger *log.Logger) {
	for range time.Tick(1 * time.Minute) {
		uniqueCountLock.Lock()

		// Log unique request count
		logger.Printf("Unique requests in the last minute: %d\n", uniqueCount)

		// Reset unique IDs and count
		uniqueIDs = make(map[int]struct{})
		uniqueCount = 0

		uniqueCountLock.Unlock()
	}
}

// fireEndpointRequest sends a GET request to the provided endpoint with unique count
func fireEndpointRequest(endpoint string, logger *log.Logger) {
	uniqueCountLock.Lock()
	currentCount := uniqueCount
	uniqueCountLock.Unlock()

	// Send HTTP GET request
	resp, err := http.Get(endpoint + "?uniqueCount=" + strconv.Itoa(currentCount))
	if err != nil {
		logger.Printf("Failed to send request to endpoint: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Log the response status
	logger.Printf("Endpoint request sent, response status: %d\n", resp.StatusCode)
}
