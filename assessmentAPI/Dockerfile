# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# This creates a statically-linked binary for better portability
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/assessmentAPI ./cmd/assessmentAPI

# Deployment stage
FROM alpine:latest

# Create a non-root user to run the application
RUN adduser -D -H -h /app appuser
USER appuser

# Set working directory
WORKDIR /app

COPY --from=builder /app/bin/assessmentAPI .
COPY .env .

EXPOSE 8080

CMD ["./assessmentAPI"]