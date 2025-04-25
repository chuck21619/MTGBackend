# Build stage
FROM golang:1.24.2 AS builder

# Set the current working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies using go mod tidy
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# Run stage
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/app /root/

# Expose the port on which your app is running
EXPOSE 8080

# Command to run the app
CMD ["/root/app"]
