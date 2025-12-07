# Dockerfile for E-Commerce Backend
# This file tells Docker how to build your application

# Stage 1: Build Stage
# We use a multi-stage build to keep the final image small
FROM golang:1.25.3-alpine AS builder

# Install build dependencies
# git: for downloading Go modules from Git repositories
# gcc, musl-dev: required for CGO (C bindings) in some packages
RUN apk add --no-cache git gcc musl-dev

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
# Docker caches layers, so if these files don't change,
# it won't re-download dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
# -o /app/main: output binary named 'main'
# -ldflags="-w -s": reduce binary size by removing debug info
# ./cmd/api: the main package location
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main -ldflags="-w -s" ./cmd/api

# Stage 2: Runtime Stage
# Use a minimal base image for the final container
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy docs folder for Swagger
COPY --from=builder /app/docs ./docs

# Create uploads directory
RUN mkdir -p uploads && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Health check to verify the container is running properly
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
