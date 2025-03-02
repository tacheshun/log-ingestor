package ingestor

import (
	"bytes"
	"encoding/json"
	"log-ingestor/internal/database"
	"log-ingestor/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *database.MockDB) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock database
	mockDB := database.NewMockDB()

	// Create a log ingestor service with the mock database
	logIngestor := NewLogIngestor(mockDB)

	// Create a router
	router := gin.Default()

	// Define routes
	router.POST("/", logIngestor.HandleLogIngestion)
	router.GET("/logs", logIngestor.QueryLogs)

	return router, mockDB
}

func TestHandleLogIngestion(t *testing.T) {
	router, mockDB := setupTestRouter()

	// Create a sample log
	timestamp, _ := time.Parse(time.RFC3339, "2023-09-15T08:00:00Z")
	log := models.Log{
		Level:      "error",
		Message:    "Failed to connect to DB",
		ResourceID: "server-1234",
		Timestamp:  timestamp,
		TraceID:    "abc-xyz-123",
		SpanID:     "span-456",
		Commit:     "5e5342f",
		Metadata: map[string]string{
			"parentResourceId": "server-0987",
		},
	}

	// Convert log to JSON
	logJSON, _ := json.Marshal(log)

	// Create a request
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(logJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "Log ingested successfully" {
		t.Errorf("Expected response status 'Log ingested successfully', got '%s'", response["status"])
	}

	// Verify the log was inserted into the mock database
	logs, _ := mockDB.QueryLogs(nil, &models.LogQuery{})
	if len(logs) != 1 {
		t.Errorf("Expected 1 log in the database, got %d", len(logs))
	}
}

func TestQueryLogs(t *testing.T) {
	router, mockDB := setupTestRouter()

	// Insert some sample logs
	timestamp1, _ := time.Parse(time.RFC3339, "2023-09-15T08:00:00Z")
	timestamp2, _ := time.Parse(time.RFC3339, "2023-09-15T09:00:00Z")

	log1 := &models.Log{
		Level:      "error",
		Message:    "Failed to connect to DB",
		ResourceID: "server-1234",
		Timestamp:  timestamp1,
		TraceID:    "abc-xyz-123",
		SpanID:     "span-456",
		Commit:     "5e5342f",
		Metadata: map[string]string{
			"parentResourceId": "server-0987",
		},
	}

	log2 := &models.Log{
		Level:      "info",
		Message:    "User authentication successful",
		ResourceID: "server-5678",
		Timestamp:  timestamp2,
		TraceID:    "def-uvw-789",
		SpanID:     "span-789",
		Commit:     "a1b2c3d",
		Metadata: map[string]string{
			"parentResourceId": "server-6543",
		},
	}

	mockDB.InsertLog(nil, log1)
	mockDB.InsertLog(nil, log2)

	// Test query with no filters
	req, _ := http.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	logs, ok := response["logs"].([]interface{})
	if !ok {
		t.Fatalf("Expected logs to be an array")
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs in the response, got %d", len(logs))
	}

	// Test query with level filter
	req, _ = http.NewRequest("GET", "/logs?level=error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	logs, ok = response["logs"].([]interface{})
	if !ok {
		t.Fatalf("Expected logs to be an array")
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log in the response, got %d", len(logs))
	}

	// Test query with resourceId filter
	req, _ = http.NewRequest("GET", "/logs?resourceId=server-5678", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	logs, ok = response["logs"].([]interface{})
	if !ok {
		t.Fatalf("Expected logs to be an array")
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log in the response, got %d", len(logs))
	}

	// Test query with message filter
	req, _ = http.NewRequest("GET", "/logs?message=Failed", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	logs, ok = response["logs"].([]interface{})
	if !ok {
		t.Fatalf("Expected logs to be an array")
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log in the response, got %d", len(logs))
	}

	// Test query with pagination
	req, _ = http.NewRequest("GET", "/logs?page=1&limit=1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	logs, ok = response["logs"].([]interface{})
	if !ok {
		t.Fatalf("Expected logs to be an array")
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log in the response, got %d", len(logs))
	}
}

// TestHandleLogIngestionWithDBError tests the error handling when the database operation fails
func TestHandleLogIngestionWithDBError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock database with error simulation
	mockDB := database.NewMockDB()
	mockDB.SimulateError = true

	// Create a log ingestor service with the mock database
	logIngestor := NewLogIngestor(mockDB)

	// Create a router
	router := gin.Default()
	router.POST("/", logIngestor.HandleLogIngestion)

	// Create a sample log
	timestamp, _ := time.Parse(time.RFC3339, "2023-09-15T08:00:00Z")
	log := models.Log{
		Level:      "error",
		Message:    "Failed to connect to DB",
		ResourceID: "server-1234",
		Timestamp:  timestamp,
		TraceID:    "abc-xyz-123",
		SpanID:     "span-456",
		Commit:     "5e5342f",
		Metadata: map[string]string{
			"parentResourceId": "server-0987",
		},
	}

	// Convert log to JSON
	logJSON, _ := json.Marshal(log)

	// Create a request
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(logJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the status code - should be an error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Failed to insert log: simulated error" {
		t.Errorf("Expected error message 'Failed to insert log: simulated error', got '%s'", response["error"])
	}
}
