# Use the official Golang image as a build stage
FROM golang:1.22.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code for task-service into the container
COPY task-service/ ./task-service

# Build the Go app
RUN go build -o task-service ./task-service

# Start a new stage from scratch
FROM gcr.io/distroless/base

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/task-service .

# Copy .env file into the container
COPY task-service/.env .env

# Command to run the executable
ENTRYPOINT ["/task-service"]
