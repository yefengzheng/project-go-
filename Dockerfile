# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o scanner ./cmd/scanner

# Runtime stage
FROM alpine:latest

# Add ca-certificates for any HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/scanner .

# Set environment variables with defaults (these can be overridden at runtime)
ENV REST_PORT=8080 \
    REST_READ_TIMEOUT=20s \
    REST_WRITE_TIMEOUT=20s \
    DOWNLOAD_WORKER_COUNT=3 \
    SCAN_WORKER_COUNT=2 \
    QUEUE_SIZE=100 \
    REDIS_ADDRESS=redis \
    REDIS_PORT=6379 \
    REDIS_PASSWORD="" \
    REDIS_DB=0 \
    MYSQL_ADDRESS=postgres \
    MYSQL_PORT=5432 \
    MYSQL_USER=postgres \
    MYSQL_PASSWORD="postgres" \
    MYSQL_CONNECT_TIMEOUT=5 \
    MYSQL_RESULT_DB=result_db

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application with the wait script
CMD ["./scanner"]
