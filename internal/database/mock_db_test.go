package database

import (
	"context"
	"log-ingestor/internal/models"
	"testing"
	"time"
)

func TestMockDBInsertAndQuery(t *testing.T) {
	// Create a new mock database
	mockDB := NewMockDB()

	// Create sample logs
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

	// Insert logs
	ctx := context.Background()
	err := mockDB.InsertLog(ctx, log1)
	if err != nil {
		t.Fatalf("Failed to insert log1: %v", err)
	}

	err = mockDB.InsertLog(ctx, log2)
	if err != nil {
		t.Fatalf("Failed to insert log2: %v", err)
	}

	// Test query with no filters
	query := &models.LogQuery{
		Page:  1,
		Limit: 10,
	}

	logs, err := mockDB.QueryLogs(ctx, query)
	if err != nil {
		t.Fatalf("Failed to query logs: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	// Test query with level filter
	levelQuery := &models.LogQuery{
		Level: "error",
		Page:  1,
		Limit: 10,
	}

	logs, err = mockDB.QueryLogs(ctx, levelQuery)
	if err != nil {
		t.Fatalf("Failed to query logs with level filter: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log with level 'error', got %d", len(logs))
	}

	if len(logs) > 0 && logs[0].Level != "error" {
		t.Errorf("Expected log with level 'error', got '%s'", logs[0].Level)
	}

	// Test query with resourceId filter
	resourceQuery := &models.LogQuery{
		ResourceID: "server-5678",
		Page:       1,
		Limit:      10,
	}

	logs, err = mockDB.QueryLogs(ctx, resourceQuery)
	if err != nil {
		t.Fatalf("Failed to query logs with resourceId filter: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log with resourceId 'server-5678', got %d", len(logs))
	}

	if len(logs) > 0 && logs[0].ResourceID != "server-5678" {
		t.Errorf("Expected log with resourceId 'server-5678', got '%s'", logs[0].ResourceID)
	}

	// Test query with time range filter
	startTime, _ := time.Parse(time.RFC3339, "2023-09-15T08:30:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2023-09-15T09:30:00Z")

	timeQuery := &models.LogQuery{
		StartTime: startTime,
		EndTime:   endTime,
		Page:      1,
		Limit:     10,
	}

	logs, err = mockDB.QueryLogs(ctx, timeQuery)
	if err != nil {
		t.Fatalf("Failed to query logs with time range filter: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log within time range, got %d", len(logs))
	}

	if len(logs) > 0 && logs[0].Timestamp.Before(startTime) {
		t.Errorf("Expected log timestamp after %v, got %v", startTime, logs[0].Timestamp)
	}

	// Test query with message filter
	messageQuery := &models.LogQuery{
		Message: "Failed",
		Page:    1,
		Limit:   10,
	}

	logs, err = mockDB.QueryLogs(ctx, messageQuery)
	if err != nil {
		t.Fatalf("Failed to query logs with message filter: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log with message containing 'Failed', got %d", len(logs))
	}

	// Test pagination
	paginationQuery := &models.LogQuery{
		Page:  1,
		Limit: 1,
	}

	logs, err = mockDB.QueryLogs(ctx, paginationQuery)
	if err != nil {
		t.Fatalf("Failed to query logs with pagination: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log with pagination, got %d", len(logs))
	}

	// Test second page
	page2Query := &models.LogQuery{
		Page:  2,
		Limit: 1,
	}

	logs, err = mockDB.QueryLogs(ctx, page2Query)
	if err != nil {
		t.Fatalf("Failed to query logs for page 2: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log on page 2, got %d", len(logs))
	}

	// Test empty page
	emptyPageQuery := &models.LogQuery{
		Page:  3,
		Limit: 1,
	}

	logs, err = mockDB.QueryLogs(ctx, emptyPageQuery)
	if err != nil {
		t.Fatalf("Failed to query logs for empty page: %v", err)
	}

	if len(logs) != 0 {
		t.Errorf("Expected 0 logs on empty page, got %d", len(logs))
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"Failed to connect to DB", "Failed", true},
		{"User authentication successful", "Failed", false},
		{"Failed", "Failed", true},
		{"", "Failed", false},
		{"Failed to connect to DB", "", true}, // Empty substring is always contained
		{"", "", true},                        // Empty substring is always contained in empty string
	}

	for _, tc := range testCases {
		t.Run(tc.s+"-"+tc.substr, func(t *testing.T) {
			result := contains(tc.s, tc.substr)
			if result != tc.expected {
				t.Errorf("contains(%q, %q) = %v, expected %v", tc.s, tc.substr, result, tc.expected)
			}
		})
	}
}
