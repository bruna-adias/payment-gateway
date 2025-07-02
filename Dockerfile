# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-gateway

# Add curl
RUN apk --no-cache add curl

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/payment-gateway .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./payment-gateway"]