package database

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log-ingestor/internal/models"
)

// MongoDB represents the MongoDB client and collection
type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Ensure MongoDB implements the DB interface
var _ DB = (*MongoDB)(nil)

// NewMongoDB creates a new MongoDB connection
func NewMongoDB() (*MongoDB, error) {
	// Get MongoDB URI from environment variable
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	// Get collection
	collection := client.Database(dbName).Collection(collectionName)

	// Create indexes for better query performance
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "level", Value: 1},
				{Key: "resourceId", Value: 1},
				{Key: "traceId", Value: 1},
				{Key: "spanId", Value: 1},
				{Key: "commit", Value: 1},
				{Key: "timestamp", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "message", Value: "text"},
			},
			Options: options.Index().SetName("message_text"),
		},
		{
			Keys: bson.D{
				{Key: "metadata.parentResourceId", Value: 1},
			},
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Printf("Error creating indexes: %v", err)
	}

	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
		return err
	}
	return nil
}

// InsertLog inserts a log into MongoDB
func (m *MongoDB) InsertLog(ctx context.Context, logEntry *models.Log) error {
	_, err := m.collection.InsertOne(ctx, logEntry)
	return err
}

// QueryLogs queries logs from MongoDB based on the provided filters
func (m *MongoDB) QueryLogs(ctx context.Context, query *models.LogQuery) ([]*models.Log, error) {
	filter := bson.M{}

	// Apply filters if provided
	if query.Level != "" {
		filter["level"] = query.Level
	}

	if query.Message != "" {
		filter["message"] = bson.M{"$regex": primitive.Regex{Pattern: query.Message, Options: "i"}}
	}

	if query.ResourceID != "" {
		filter["resourceId"] = query.ResourceID
	}

	if query.TraceID != "" {
		filter["traceId"] = query.TraceID
	}

	if query.SpanID != "" {
		filter["spanId"] = query.SpanID
	}

	if query.Commit != "" {
		filter["commit"] = query.Commit
	}

	if query.ParentResourceID != "" {
		filter["metadata.parentResourceId"] = query.ParentResourceID
	}

	// Date range filter
	timeFilter := bson.M{}
	if !query.StartTime.IsZero() {
		timeFilter["$gte"] = query.StartTime
	}
	if !query.EndTime.IsZero() {
		timeFilter["$lte"] = query.EndTime
	}
	if len(timeFilter) > 0 {
		filter["timestamp"] = timeFilter
	}

	// Regex pattern search
	if query.RegexPattern != "" {
		// Validate regex pattern
		_, err := regexp.Compile(query.RegexPattern)
		if err == nil {
			filter["message"] = bson.M{"$regex": primitive.Regex{Pattern: query.RegexPattern, Options: "i"}}
		}
	}

	// Full-text search
	if query.FullTextSearch != "" {
		filter["$text"] = bson.M{"$search": query.FullTextSearch}
	}

	// Set default pagination values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	// Calculate skip value for pagination
	skip := (query.Page - 1) * query.Limit

	// Set options
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.Limit)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	// Execute query
	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var logs []*models.Log
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
