# Stage 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download necessary Go modules
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Runtime Environment
FROM debian:bullseye-slim

# Install necessary packages for the runtime
RUN apt-get update && apt-get install -y \
    libasound2 \
    ffmpeg \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install edge-tts using npm
RUN apt-get install -y npm && npm install -g edge-tts

# Set the working directory for the runtime
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/main /app/

# Expose the service port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
