# Build stage
FROM golang:1.22.4 AS builder
WORKDIR /app
COPY . .
RUN go build -o app

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
