package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// LogEntry represents the structure of a log entry
type LogEntry struct {
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	ResourceID string            `json:"resourceId"`
	Timestamp  time.Time         `json:"timestamp"`
	TraceID    string            `json:"traceId"`
	SpanID     string            `json:"spanId"`
	Commit     string            `json:"commit"`
	Metadata   map[string]string `json:"metadata"`
}

var (
	levels = []string{"error", "warning", "info", "debug"}

	messages = []string{
		"Failed to connect to DB",
		"API request timed out",
		"User authentication successful",
		"Cache miss for key",
		"Background job completed",
		"Service started successfully",
		"Memory usage high",
		"CPU usage exceeded threshold",
		"File not found",
		"Network connection lost",
	}

	resourceIDs = []string{
		"server-1234",
		"server-5678",
		"api-gateway-1",
		"worker-process-3",
		"cache-node-2",
	}

	commits = []string{
		"5e5342f",
		"a1b2c3d",
		"f4e5d6c",
		"1a2b3c4",
		"9z8y7x6",
	}

	parentResourceIDs = []string{
		"server-0987",
		"server-6543",
		"api-cluster-1",
		"worker-pool-2",
		"cache-cluster-3",
	}
)

func generateRandomLog() LogEntry {
	now := time.Now()
	startTime := now.Add(-24 * time.Hour)
	randomTime := startTime.Add(time.Duration(rand.Int63n(int64(now.Sub(startTime)))))

	traceID := fmt.Sprintf("trace-%s", randomString(6))
	spanID := fmt.Sprintf("span-%s", randomString(3))

	return LogEntry{
		Level:      levels[rand.Intn(len(levels))],
		Message:    messages[rand.Intn(len(messages))],
		ResourceID: resourceIDs[rand.Intn(len(resourceIDs))],
		Timestamp:  randomTime,
		TraceID:    traceID,
		SpanID:     spanID,
		Commit:     commits[rand.Intn(len(commits))],
		Metadata: map[string]string{
			"parentResourceId": parentResourceIDs[rand.Intn(len(parentResourceIDs))],
		},
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	// No need to seed the random number generator in Go 1.20+
	// The math/rand package is automatically seeded now

	// Number of logs to generate
	numLogs := 100

	// Log ingestor endpoint
	endpoint := "http://localhost:3000"

	fmt.Printf("Generating and sending %d logs to %s...\n", numLogs, endpoint)

	// Generate and send logs
	for i := 0; i < numLogs; i++ {
		logEntry := generateRandomLog()

		// Convert log to JSON
		logJSON, err := json.Marshal(logEntry)
		if err != nil {
			log.Fatalf("Error marshaling log: %v", err)
		}

		// Send log to ingestor
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(logJSON))
		if err != nil {
			log.Fatalf("Error sending log: %v", err)
		}

		// Check response
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Error response: %v", resp.Status)
		}

		resp.Body.Close()

		// Print progress
		if (i+1)%10 == 0 {
			fmt.Printf("Sent %d logs...\n", i+1)
		}

		// Add a small delay to avoid overwhelming the server
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Done! All logs sent successfully.")
}
