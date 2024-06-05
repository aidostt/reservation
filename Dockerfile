# First stage: build the application
FROM golang:1.20-alpine AS builder

# Install necessary packages
RUN apk update && apk add --no-cache git ca-certificates wget && update-ca-certificates

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o api ./cmd/api

# Download and install the migration tool
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz \
    && tar -xzf migrate.linux-amd64.tar.gz \
    && mv migrate.linux-amd64 /usr/local/bin/migrate \
    && rm migrate.linux-amd64.tar.gz

# Second stage: run the application
FROM alpine:latest

# Install necessary certificates
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Set the current working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Expose the necessary port
EXPOSE 9090

# Command to run the migration and the application
CMD ["sh", "-c", "migrate -path /app/migrations -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable up && ./api"]
