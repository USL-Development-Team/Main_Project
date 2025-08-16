# Multi-stage build for USL Go application
FROM node:18-alpine AS tailwind-builder

# Install Tailwind CSS and build assets
WORKDIR /app
COPY package*.json ./
RUN npm install

# Copy source files needed for Tailwind build
COPY static/ ./static/
COPY templates/ ./templates/
COPY tailwind.config.js ./

# Build Tailwind CSS
RUN npm run build

# Go build stage
FROM golang:1.24.2-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built CSS from tailwind stage
COPY --from=tailwind-builder /app/static/dist/ ./static/dist/

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w -extldflags "-static"' -o server ./cmd/server

# Final runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=go-builder /app/server .

# Copy static assets and templates
COPY --from=go-builder /app/static/ ./static/
COPY --from=go-builder /app/templates/ ./templates/

# Create logs directory
RUN mkdir -p logs && chown app:app logs

# Change ownership to app user
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port (Render will provide PORT env var)
EXPOSE 8080

# Health check - use dedicated health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider --timeout=2 http://localhost:${PORT:-8080}/health || exit 1

# Run the binary
CMD ["./server"]