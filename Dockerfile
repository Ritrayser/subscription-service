# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Добавляем строку ниже
RUN go mod tidy
RUN go build -o /app/server ./cmd/server

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY .env.example .env
EXPOSE 8080
CMD ["/app/server"]