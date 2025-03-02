package models

import (
	"time"
)

// Log represents the structure of a log entry
type Log struct {
	Level      string            `json:"level" bson:"level"`
	Message    string            `json:"message" bson:"message"`
	ResourceID string            `json:"resourceId" bson:"resourceId"`
	Timestamp  time.Time         `json:"timestamp" bson:"timestamp"`
	TraceID    string            `json:"traceId" bson:"traceId"`
	SpanID     string            `json:"spanId" bson:"spanId"`
	Commit     string            `json:"commit" bson:"commit"`
	Metadata   map[string]string `json:"metadata" bson:"metadata"`
}

// LogQuery represents the query parameters for filtering logs
type LogQuery struct {
	Level            string    `form:"level"`
	Message          string    `form:"message"`
	ResourceID       string    `form:"resourceId"`
	TraceID          string    `form:"traceId"`
	SpanID           string    `form:"spanId"`
	Commit           string    `form:"commit"`
	ParentResourceID string    `form:"parentResourceId"`
	StartTime        time.Time `form:"startTime" time_format:"2006-01-02T15:04:05Z"`
	EndTime          time.Time `form:"endTime" time_format:"2006-01-02T15:04:05Z"`
	RegexPattern     string    `form:"regex"`
	FullTextSearch   string    `form:"search"`
	Page             int       `form:"page"`
	Limit            int       `form:"limit"`
}
