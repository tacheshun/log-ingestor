package ingestor

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"log-ingestor/internal/database"
	"log-ingestor/internal/models"
)

// LogIngestor represents the log ingestor service
type LogIngestor struct {
	db database.DB
}

// NewLogIngestor creates a new log ingestor service
func NewLogIngestor(db database.DB) *LogIngestor {
	return &LogIngestor{
		db: db,
	}
}

// HandleLogIngestion handles the log ingestion HTTP request
func (li *LogIngestor) HandleLogIngestion(c *gin.Context) {
	var logEntry models.Log

	// Bind JSON request body to log entry
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure timestamp is valid
	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now().UTC()
	}

	// Insert log into database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := li.db.InsertLog(ctx, &logEntry); err != nil {
		log.Printf("Error inserting log: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Log ingested successfully"})
}

// QueryLogs handles the log query HTTP request
func (li *LogIngestor) QueryLogs(c *gin.Context) {
	var query models.LogQuery

	// Bind query parameters to log query
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query logs from database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logs, err := li.db.QueryLogs(ctx, &query)
	if err != nil {
		log.Printf("Error querying logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "count": len(logs)})
}
