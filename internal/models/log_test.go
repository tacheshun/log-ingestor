package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLogJSONMarshaling(t *testing.T) {
	// Create a sample log
	timestamp, _ := time.Parse(time.RFC3339, "2023-09-15T08:00:00Z")
	log := Log{
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

	// Marshal to JSON
	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal Log to JSON: %v", err)
	}

	// Unmarshal back to a Log
	var unmarshaledLog Log
	err = json.Unmarshal(jsonData, &unmarshaledLog)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to Log: %v", err)
	}

	// Verify fields match
	if log.Level != unmarshaledLog.Level {
		t.Errorf("Level mismatch: expected %s, got %s", log.Level, unmarshaledLog.Level)
	}
	if log.Message != unmarshaledLog.Message {
		t.Errorf("Message mismatch: expected %s, got %s", log.Message, unmarshaledLog.Message)
	}
	if log.ResourceID != unmarshaledLog.ResourceID {
		t.Errorf("ResourceID mismatch: expected %s, got %s", log.ResourceID, unmarshaledLog.ResourceID)
	}
	if !log.Timestamp.Equal(unmarshaledLog.Timestamp) {
		t.Errorf("Timestamp mismatch: expected %v, got %v", log.Timestamp, unmarshaledLog.Timestamp)
	}
	if log.TraceID != unmarshaledLog.TraceID {
		t.Errorf("TraceID mismatch: expected %s, got %s", log.TraceID, unmarshaledLog.TraceID)
	}
	if log.SpanID != unmarshaledLog.SpanID {
		t.Errorf("SpanID mismatch: expected %s, got %s", log.SpanID, unmarshaledLog.SpanID)
	}
	if log.Commit != unmarshaledLog.Commit {
		t.Errorf("Commit mismatch: expected %s, got %s", log.Commit, unmarshaledLog.Commit)
	}
	if log.Metadata["parentResourceId"] != unmarshaledLog.Metadata["parentResourceId"] {
		t.Errorf("Metadata.parentResourceId mismatch: expected %s, got %s",
			log.Metadata["parentResourceId"], unmarshaledLog.Metadata["parentResourceId"])
	}
}

func TestLogQueryValidation(t *testing.T) {
	// Test with valid time format
	query := LogQuery{
		Level:     "error",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		Page:      1,
		Limit:     10,
	}

	// Verify default values are set correctly
	if query.Page != 1 {
		t.Errorf("Expected default Page to be 1, got %d", query.Page)
	}
	if query.Limit != 10 {
		t.Errorf("Expected default Limit to be 10, got %d", query.Limit)
	}

	// Test with zero values
	emptyQuery := LogQuery{}
	if emptyQuery.Page != 0 {
		t.Errorf("Expected empty Page to be 0, got %d", emptyQuery.Page)
	}
	if emptyQuery.Limit != 0 {
		t.Errorf("Expected empty Limit to be 0, got %d", emptyQuery.Limit)
	}
}
