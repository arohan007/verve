package utils

import (
	"log"
	"os"
)

// InitLogger initializes a logger to write logs to a file
func InitLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	return log.New(logFile, "", log.LstdFlags)
}

