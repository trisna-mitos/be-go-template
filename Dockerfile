# Development stage
FROM golang:1.24.1-alpine AS development

# Install necessary packages
RUN apk add --no-cache git

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose ports
EXPOSE 8080 50051

# Start Air for hot reload
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files for Swagger UI
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./docs

# Expose ports
EXPOSE 8080 50051

# Run the binary
CMD ["./main"]