# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Build the application
RUN go mod download
RUN go build -o telegame-backend .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and config from builder
COPY --from=builder /app/telegame-backend .
COPY --from=builder /app/example.config.yaml ./config.yaml

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./telegame-backend"]
