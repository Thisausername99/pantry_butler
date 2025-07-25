# Build stage
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . /app
RUN go build -o main ./cmd/pantry_butler/server.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 5000
CMD ["/app/main"]
