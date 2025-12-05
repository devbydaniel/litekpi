# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy migrations if they exist
COPY --from=builder /app/migrations ./migrations

# Change ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
