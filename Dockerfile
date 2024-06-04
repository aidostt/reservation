# First stage: build the application
FROM golang:1.20 AS builder

# Install necessary packages
RUN apt-get update && apt-get install -y git ca-certificates wget tar

# Install migration tool
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz && \
    mkdir -p /tmp/migrate && \
    tar -xzf migrate.linux-amd64.tar.gz -C /tmp/migrate && \
    mv /tmp/migrate/migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm -rf /tmp/migrate migrate.linux-amd64.tar.gz

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o api ./cmd/api

# Second stage: run the application
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache ca-certificates

# Set the current working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Copy the migrate binary from the builder stage
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Ensure the migrate binary is executable
RUN chmod +x /usr/local/bin/migrate

# Expose the necessary port
EXPOSE 9090

# Command to run migrations and start the application
CMD /usr/local/bin/migrate -path ./migrations -database "${POSTGRES_URI}" up && ./api
