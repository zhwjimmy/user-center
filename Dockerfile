# Build stage
FROM golang:1.23.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached if go.mod/go.sum don't change)
RUN go mod download

# Install build tools (cached layer)
RUN go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code (only necessary files)
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY configs/ ./configs/
COPY migrations/ ./migrations/
COPY locales/ ./locales/

# Generate wire code and swagger docs
RUN wire ./cmd/usercenter && \
    swag init -g cmd/usercenter/main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o usercenter ./cmd/usercenter

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S usercenter && \
    adduser -u 1001 -S usercenter -G usercenter

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/usercenter .

# Copy configuration files
COPY --from=builder /app/configs ./configs

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Copy localization files
COPY --from=builder /app/locales ./locales

# Create logs directory
RUN mkdir -p logs && chown -R usercenter:usercenter /app

# Switch to non-root user
USER usercenter

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./usercenter"] 