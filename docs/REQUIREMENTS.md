# Log Ingestor and Query Interface


## Objective

Develop a log ingestor system that can efficiently handle efficiently vast volumes of log data and offer a simple interface for querying the data using full-text search or specific field filters.
The log ingestor must be built using Go programming language. The Query interface can by only HTML/CSS/JavaScript or React app.

The logs should be ingested (in the log ingestor) over HTTP, on port `3000`.

> We need a script to populate the logs into your system, so please ensure that the default 
> port is set to the port mentioned above as environment variable.


### Sample Log Data Format:

The logs to be ingested will be sent in this format.

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

## Requirements

The requirements are to code a log ingestor and the query interface are specified below.

### Log Ingestor:

- Develop a mechanism to ingest logs in the provided format.
- Ensure scalability to handle high volumes of logs efficiently.
- Mitigate potential bottlenecks such as I/O operations, database write speeds, etc.
- Make sure that the logs are ingested via an HTTP server, which runs on port `3000` by default.

### Query Interface:

- Offer a user interface (Web UI or CLI) for full-text search across logs.
- Include filters based on:
    - level
    - message
    - resourceId
    - traceId
    - spanId
    - commit
    - metadata.parentResourceId
- Aim for efficient and quick search results.

## Advanced Features:
- Implement search within specific date ranges.
- Utilize regular expressions for search.
- Allow combining multiple filters at the same time. For example, filtering for level=error AND resourceId=server-1234
- Provide real-time log ingestion and searching capabilities.
- Implement role-based access to the query interface.

## Sample Queries

The following are some sample queries that can be used as examples for testing

- Find all logs with the level set to "error".
- Search for logs with the message containing the term "Failed to connect".
- Retrieve all logs related to resourceId "server-1234".
- Filter logs between the timestamp "2023-09-10T00:00:00Z" and "2023-09-15T23:59:59Z". (Advanced Feature)

## Evaluation Criteria:

- Volume: The ability of your system to ingest massive volumes of data.
- Speed: Efficiency in returning search results.
- Scalability: Adaptability to increasing volumes of logs/queries.
- Readability: The cleanliness and structure of the codebase.
- Usability: Intuitive, user-friendly interface.
- Advanced Features: Implementation of bonus functionalities.

## Tips:
Here are a few tips for completing the specified task.
- Consider hybrid database solutions (relational + NoSQL) for a balance of structured data handling and efficient search capabilities.
- Database indexing and sharding might be beneficial for scalability and speed.
- Distributed systems or cloud-based solutions can ensure robust scalability.