# Use official Go image as a base
FROM golang:1.23

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the rest of your backend code
COPY . .

# Install Go (if not already included in the base image)
# RUN go install ./scripts/seed-user

# Build nothing yet — we'll run things at container start
CMD ["go", "run", "main.go"]
