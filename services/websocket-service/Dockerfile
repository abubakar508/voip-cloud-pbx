# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/websocket-service ./...

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/websocket-service /app/websocket-service

ENV WEBSOCKET_SERVICE_PORT=8084

EXPOSE 8084

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8084/healthz" || exit 1

ENTRYPOINT ["/app/websocket-service"]
