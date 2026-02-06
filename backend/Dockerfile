# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o app .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"]
