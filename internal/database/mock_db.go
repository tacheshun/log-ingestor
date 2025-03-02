package database

import (
	"context"
	"errors"
	"log-ingestor/internal/models"
	"strings"
	"sync"
)

// MockDB is a mock implementation of the DB interface for testing
type MockDB struct {
	logs          []*models.Log
	mutex         sync.RWMutex
	SimulateError bool
}

// Ensure MockDB implements the DB interface
var _ DB = (*MockDB)(nil)

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		logs:          make([]*models.Log, 0),
		SimulateError: false,
	}
}

// Close is a no-op for the mock database
func (m *MockDB) Close() error {
	// No-op
	return nil
}

// InsertLog inserts a log into the mock database
func (m *MockDB) InsertLog(ctx context.Context, logEntry *models.Log) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.SimulateError {
		return errors.New("simulated error")
	}

	m.logs = append(m.logs, logEntry)
	return nil
}

// QueryLogs queries logs from the mock database based on the provided filters
func (m *MockDB) QueryLogs(ctx context.Context, query *models.LogQuery) ([]*models.Log, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.SimulateError {
		return nil, errors.New("simulated error")
	}

	// Apply default pagination values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	// Filter logs based on query parameters
	var filteredLogs []*models.Log
	for _, log := range m.logs {
		if matchesQuery(log, query) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// Apply pagination
	start := (query.Page - 1) * query.Limit
	end := start + query.Limit

	if start >= len(filteredLogs) {
		return []*models.Log{}, nil
	}

	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}

	return filteredLogs[start:end], nil
}

// matchesQuery checks if a log matches the query parameters
func matchesQuery(log *models.Log, query *models.LogQuery) bool {
	// Level filter
	if query.Level != "" && log.Level != query.Level {
		return false
	}

	// ResourceID filter
	if query.ResourceID != "" && log.ResourceID != query.ResourceID {
		return false
	}

	// TraceID filter
	if query.TraceID != "" && log.TraceID != query.TraceID {
		return false
	}

	// SpanID filter
	if query.SpanID != "" && log.SpanID != query.SpanID {
		return false
	}

	// Commit filter
	if query.Commit != "" && log.Commit != query.Commit {
		return false
	}

	// ParentResourceID filter
	if query.ParentResourceID != "" && log.Metadata["parentResourceId"] != query.ParentResourceID {
		return false
	}

	// Date range filter
	if !query.StartTime.IsZero() && log.Timestamp.Before(query.StartTime) {
		return false
	}
	if !query.EndTime.IsZero() && log.Timestamp.After(query.EndTime) {
		return false
	}

	// Simple message contains filter (not regex for simplicity)
	if query.Message != "" {
		// Simple contains check for testing
		if !contains(log.Message, query.Message) {
			return false
		}
	}

	return true
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	if substr == "" {
		return true // Empty substring is always contained
	}

	return strings.Contains(s, substr)
}
