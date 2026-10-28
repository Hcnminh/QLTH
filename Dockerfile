# Multi-stage build for Go Core
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# Final lightweight scratch/alpine image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app
COPY --from=builder /app/bin/api /app/api

EXPOSE 8080

CMD ["/app/api"]
