# ==============================================================================
# Build stage
# ==============================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/...

# ==============================================================================
# Production stage
# ==============================================================================
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS requests to external APIs
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /server /app/server

# Create non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run the binary
CMD ["/app/server"]
