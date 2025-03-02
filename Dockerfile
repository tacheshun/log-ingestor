# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o log-ingestor .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/log-ingestor .

# Copy UI files
COPY --from=builder /app/ui/dist ./ui/dist

# Expose the port
EXPOSE 3000

# Set environment variables
ENV PORT=3000
ENV GIN_MODE=release

# Run the application
CMD ["./log-ingestor"] 