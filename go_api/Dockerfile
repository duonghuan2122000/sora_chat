# syntax=docker/dockerfile:1

# Stage 1: build app
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cài đặt git để go mod tidy chạy được
RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/server

# Stage 2: tạo image nhỏ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
# Copy file config.yaml từ source code
COPY --from=builder /app/config.yaml .
EXPOSE 8080

CMD ["./main"]