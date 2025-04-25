# Build stage
FROM golang:1.24.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the codebase (both backend code and frontend files)
COPY . .

# Build the Go binary from the backend directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./backend

# Run stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/app .

# Copy the frontend files into the container so Go can serve them
COPY --from=builder /app/frontend ./frontend

# Expose the port your app listens on
EXPOSE 8080

# Start the application
CMD ["./app"]
