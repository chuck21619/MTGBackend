# Build stage
FROM golang:1.22.4 AS builder
WORKDIR /app
COPY . .

# 🔧 Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app /root/
EXPOSE 8080
CMD ["/root/app"]
