# Build stage
FROM golang:1.24.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate wire code
RUN go install github.com/google/wire/cmd/wire@latest
RUN wire ./cmd/usercenter

# Generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/usercenter/main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o usercenter ./cmd/usercenter

# Final stage
FROM alpine:3.22

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