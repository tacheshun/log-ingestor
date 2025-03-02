# Log Ingestor and Query Interface

A system for efficiently ingesting and querying log data with a user-friendly interface.

## Features

- HTTP-based log ingestion on port 3000
- MongoDB for efficient storage and retrieval of logs
- Full-text search capabilities
- Advanced filtering options:
  - Level (error, warning, info, debug)
  - Message content
  - Resource ID
  - Trace ID
  - Span ID
  - Commit
  - Parent Resource ID
- Date range filtering
- Regular expression search
- Real-time log ingestion
- Responsive web UI for querying logs

## Requirements

- Go 1.18 or higher
- MongoDB
- Docker and Docker Compose (optional, for running MongoDB)

## Setup and Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd log-ingestor
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Start MongoDB using Docker Compose:
   ```
   docker-compose up -d
   ```

4. Run the application:
   ```
   go run main.go
   ```

5. The server will start on port 3000 by default. You can access the query interface at:
   ```
   http://localhost:3000
   ```

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
PORT=3000
MONGODB_URI=mongodb://localhost:27017
DB_NAME=log_ingestor
COLLECTION_NAME=logs
```

## API Endpoints

### Log Ingestion

- **URL**: `/`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "level": "error",
    "message": "Failed to connect to DB",
    "resourceId": "server-1234",
    "timestamp": "2023-09-15T08:00:00Z",
    "traceId": "abc-xyz-123",
    "spanId": "span-456",
    "commit": "5e5342f",
    "metadata": {
      "parentResourceId": "server-0987"
    }
  }
  ```

### Query Logs

- **URL**: `/logs`
- **Method**: `GET`
- **Query Parameters**:
  - `level`: Filter by log level
  - `message`: Filter by message content
  - `resourceId`: Filter by resource ID
  - `traceId`: Filter by trace ID
  - `spanId`: Filter by span ID
  - `commit`: Filter by commit
  - `parentResourceId`: Filter by parent resource ID
  - `startTime`: Filter logs after this time (ISO format)
  - `endTime`: Filter logs before this time (ISO format)
  - `regex`: Search using regular expression
  - `search`: Full-text search
  - `page`: Page number for pagination
  - `limit`: Number of logs per page

## Sample Queries

1. Find all logs with the level set to "error":
   ```
   GET /logs?level=error
   ```

2. Search for logs with the message containing the term "Failed to connect":
   ```
   GET /logs?message=Failed%20to%20connect
   ```

3. Retrieve all logs related to resourceId "server-1234":
   ```
   GET /logs?resourceId=server-1234
   ```

4. Filter logs between a specific time range:
   ```
   GET /logs?startTime=2023-09-10T00:00:00Z&endTime=2023-09-15T23:59:59Z
   ```

## Architecture

The application follows a clean architecture pattern:

- **Models**: Define the data structures for logs and queries
- **Database**: Handles MongoDB connection and operations
- **Ingestor**: Manages HTTP request handling for log ingestion and querying
- **UI**: Provides a user-friendly interface for querying logs

## Performance Considerations

- MongoDB indexes are created for efficient querying
- Pagination is implemented to handle large result sets
- The UI is designed to be responsive and efficient
- The server uses Goroutines for concurrent request handling

## License

[MIT License](LICENSE)