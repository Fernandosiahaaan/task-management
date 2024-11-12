# Use the official Golang image as a build stage
FROM golang:1.22.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
COPY ./ ./
RUN go build -o user-service .
FROM gcr.io/distroless/base
WORKDIR /
COPY --from=builder /app/user-service ./
COPY .env .env

# Define the entrypoint for the application
ENTRYPOINT ["/user-service"]
