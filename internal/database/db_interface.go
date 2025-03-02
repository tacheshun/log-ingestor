package database

import (
	"context"
	"log-ingestor/internal/models"
)

// DB is an interface for database operations
type DB interface {
	// Close closes the database connection
	Close() error

	// InsertLog inserts a log into the database
	InsertLog(ctx context.Context, logEntry *models.Log) error

	// QueryLogs queries logs from the database based on the provided filters
	QueryLogs(ctx context.Context, query *models.LogQuery) ([]*models.Log, error)
}
