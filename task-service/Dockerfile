# Use the official Golang image as a build stage
FROM golang:1.22.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o task-service .
FROM gcr.io/distroless/base
COPY --from=builder /app/task-service ./
COPY .env .env

# Command to run the executable
ENTRYPOINT ["/task-service"]
