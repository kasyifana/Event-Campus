# Multi-stage build for Go application

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download && go mod verify

# Install swag for swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate swagger docs
RUN /go/bin/swag init -g cmd/api/main.go --output docs

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.version=${BUILD_VERSION:-dev}" \
    -o main cmd/api/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Set timezone to Asia/Jakarta
ENV TZ=Asia/Jakarta

WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Create storage directories with proper permissions
RUN mkdir -p storage/posters storage/documents && \
    chown -R appuser:appuser storage

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations (if needed for auto-migration)
COPY --from=builder /app/migrations ./migrations

# Change to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]
