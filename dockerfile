# Build stage
FROM golang:1.24.2@sha256:30baaea08c5d1e858329c50f29fe381e9b7d7bced11a0f5f1f69a1504cdfbf5e AS builder

WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build both migrate and app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install pg_isready (comes with postgres client tools)
RUN apk add --no-cache postgresql-client

# Copy the built binaries from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/migrate .

# Copy the wait script
COPY wait-for-postgres.sh .

# Make it executable
RUN chmod +x wait-for-postgres.sh

# Expose the application port
EXPOSE 8080

# Run the wait script first, then run migrate and app
CMD ["./wait-for-postgres.sh", "sh", "-c", "./migrate && ./app"]
