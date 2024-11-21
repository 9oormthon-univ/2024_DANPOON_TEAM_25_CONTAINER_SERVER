FROM golang:1.22-alpine AS builder

# Install necessary dependencies
RUN apk add --no-cache git docker-cli

WORKDIR /app

# Copy and download dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/grpc

# Final stage
FROM alpine:latest

# Install docker-cli for docker socket interactions
RUN apk add --no-cache docker-cli

# Copy the built binary from the builder stage
COPY --from=builder /app/server /app/server

# Create an entrypoint script to handle potential docker socket permissions
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'if [ -S /var/run/docker.sock ]; then' >> /entrypoint.sh && \
    echo '  chmod 666 /var/run/docker.sock' >> /entrypoint.sh && \
    echo 'fi' >> /entrypoint.sh && \
    echo '/app/server' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Expose the gRPC port
EXPOSE 50051

# Set the entrypoint
CMD ["/bin/sh", "/entrypoint.sh"]
