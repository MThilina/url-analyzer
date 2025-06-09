# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install Git for Go modules (if needed)
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o url-analyzer ./cmd

# Stage 2: Run
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy necessary files from builder
COPY --from=builder /app/url-analyzer .
COPY templates/ ./templates/
COPY static/ ./static/
COPY config/ ./config/
COPY docs/ ./docs/

# Port from config.yaml — typically 8080
EXPOSE 8080

# Run the application
ENTRYPOINT ["./url-analyzer"]
